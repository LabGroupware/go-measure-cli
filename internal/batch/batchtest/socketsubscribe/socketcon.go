package socketsubscribe

import (
	"context"
	"encoding/csv"
	"fmt"
	"sync"
	"time"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
	"github.com/LabGroupware/go-measure-tui/internal/ws"
	"github.com/jmespath/go-jmespath"
)

func SocketConnect(
	ctx context.Context,
	ctr *app.Container,
	conf SocketConnectConfig,
	store *sync.Map,
	outputRoot string,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer fmt.Println("Socket Closed------------------------------------------------------------")

	dataTermChan := make(chan DataTypeChan)
	var breakTime time.Duration

	if conf.Term.Time != nil {
		var err error
		breakTime, err = time.ParseDuration(*conf.Term.Time)
		if err != nil {
			return fmt.Errorf("failed to parse time: %v", err)
		}
	}

	socketStart := time.Now()

	msgHandler := func(s *ws.WebSocket, msg *ws.EventResponseMessage, raw []byte) error {

		sock, err := GlobalSock.FindSocket(conf.ID)
		if err != nil {
			ctr.Logger.Error(ctx, "failed to find socket",
				logger.Value("error", err))
			return fmt.Errorf("failed to find socket: %v", err)
		}

		var selfConsumer string
		for _, consumer := range sock.Consumers {
			if v, ok := sock.ConsumerSelfEventFilterMap[consumer]; ok {
				if v.JMESPath != nil {
					if result, err := v.JMESPath.Search(msg.Data); err == nil && result != nil {
						res, ok := result.(bool)
						if ok && res {
							selfConsumer = consumer
							break
						}
					}
				}
			}
		}

		for _, action := range sock.GetActions() {

			if actionConsumer, ok := sock.GetConsumerIDByActionID(action.ID); !ok || actionConsumer != selfConsumer {
				continue
			}

			if !ws.ContainsEventType(action.EventTypes, msg.EventType) {
				continue
			}

			dataMap := make(map[string]string)

			for _, dataConf := range action.Data {
				var result any
				var err error

				if dataConf.JMESPath != "" {
					jmesPathQuery := dataConf.JMESPath
					result, err = jmespath.Search(jmesPathQuery, msg.Data)
					if err != nil {
						ctr.Logger.Error(ctx, "failed to search jmespath",
							logger.Value("error", err), logger.Value("on", "SocketSubscribe"))
						switch dataConf.OnError {
						case "error":
							dataTermChan <- DataChanError
							return nil
						}
					}
					if result == nil {
						switch dataConf.OnNil {
						case "error":
							dataTermChan <- DataChanError
							return nil
						}
						result = ""
					}
				} else {
					switch dataConf.OnError {
					case "error":
						dataTermChan <- DataChanError
						return nil
					}
					result = ""
				}

				dataMap[dataConf.Key] = fmt.Sprintf("%v", result)
			}

			if ContainsSocketActionType(action.Types, SocketActionTypeOutput) {
				file, ok := sock.GetActionsFileMap(action.ID)
				if !ok {
					ctr.Logger.Error(ctx, "failed to get file",
						logger.Value("actionID", action.ID))
					dataTermChan <- DataChanError
					return nil
				}
				// fmt.Println("Writing to file", file.Name())
				writer := csv.NewWriter(file)
				data := []string{
					msg.EventType.String(),
					socketStart.Format(time.RFC3339),
					time.Now().Format(time.RFC3339),
					fmt.Sprintf("%dms", time.Since(socketStart).Milliseconds()),
				}

				for _, dataConf := range action.Data {
					data = append(data, dataMap[dataConf.Key])
				}

				if err := writer.Write(data); err != nil {
					dataTermChan <- DataChanError
					return nil
				}
				writer.Flush()
			}

			if ContainsSocketActionType(action.Types, SocketActionTypeStore) {
				for k, v := range dataMap {
					ctr.Logger.Info(ctx, "store data",
						logger.Value("key", k), logger.Value("value", v))
					store.Store(k, v)
				}
			}

			if ContainsSocketActionType(action.Types, SocketActionTypeUnsubscribe) {
				if err := sock.UnsubscribeNotifyByAction(ctx, action.ID); err != nil {
					ctr.Logger.Error(ctx, "failed to unsubscribe notify",
						logger.Value("error", err))
					return fmt.Errorf("failed to unsubscribe notify: %v", err)
				}
			}
		}

		for _, event := range conf.Term.Event {
			if ws.ContainsEventType(event.Types, msg.EventType) {
				var termType DataTypeChan
				if event.Success {
					termType = DataTypeChanSuccessEvent
				} else {
					termType = DataTypeChanFailEvent
				}

				if event.JMESPath != "" {
					jmesPathQuery := event.JMESPath
					result, err := jmespath.Search(jmesPathQuery, msg.Data)
					if err != nil {
						ctr.Logger.Error(ctx, "failed to search jmespath",
							logger.Value("error", err), logger.Value("on", "SocketSubscribe"))
					}
					if result != nil {
						if v, ok := result.(bool); ok {
							if v {
								ctr.Logger.Info(ctx, "jmespath query result is true",
									logger.Value("on", "SocketSubscribe"))
								dataTermChan <- termType
								return nil
							}
						} else {
							ctr.Logger.Warn(ctx, "The result of the jmespath query is not a boolean",
								logger.Value("on", "SocketSubscribe"))
						}
					}
				} else {
					dataTermChan <- termType
					return nil
				}
			}
		}
		return nil
	}

	subscribeHandler := func(ws *ws.WebSocket, msg *ws.SubscribeResponseMessage) error {
		ctr.Logger.Debug(ctx, "subscribed",
			logger.Value("subscription_id", msg.SubscriptionID))

		return nil
	}

	sock, err := NewSocketWithSubscribeHandler(ctx, ctr, msgHandler, subscribeHandler)
	if err != nil {
		ctr.Logger.Error(ctx, "failed to create socket",
			logger.Value("error", err))
		return fmt.Errorf("failed to create socket: %v", err)
	}

	done, err := sock.Connect(ctx, ctr, ws.ConnectConfig{
		DisconnectOnReadMsgError:       ContainsTermError(conf.Term.Error, ErrorTypeForTermReadError),
		DisconnectOnUnmarshalJSONError: ContainsTermError(conf.Term.Error, ErrorTypeForTermUnmarshalError),
		DisconnectOnParseMsgError:      ContainsTermError(conf.Term.Error, ErrorTypeForTermParseError),
	})
	if err != nil {
		ctr.Logger.Error(ctx, "failed to connect to socket",
			logger.Value("error", err))
		return fmt.Errorf("failed to connect to socket: %v", err)
	}

	GlobalSock.AddSocket(
		conf.ID,
		NewSock(sock, conf.Output.Enabled),
	)
	defer GlobalSock.CloseSocket(conf.ID)

	var timeout <-chan time.Time
	if breakTime > 0 {
		timeout = time.After(breakTime)
	}

	select {
	case termType := <-done:
		switch termType {
		case ws.TerminateTypeConnectionClosed:
			if ContainsSuccessTerm(conf.SuccessTerm, SuccessTermClose) {
				ctr.Logger.Info(ctx, "socket connection closed")
				return nil
			}
			ctr.Logger.Warn(ctx, "socket connection closed")
			return fmt.Errorf("socket connection closed")
		default:
			if ContainsSuccessTerm(conf.SuccessTerm, SuccessTermError) {
				ctr.Logger.Warn(ctx, "socket connection terminated",
					logger.Value("type", termType))
				return nil
			}
			ctr.Logger.Warn(ctx, "socket connection terminated",
				logger.Value("type", termType))
			return fmt.Errorf("socket connection terminated")
		}
	case dataType := <-dataTermChan:
		switch dataType {
		case DataTypeChanSuccessEvent:
			ctr.Logger.Info(ctx, "success event received")
			return nil
		case DataTypeChanFailEvent:
			ctr.Logger.Warn(ctx, "fail event received")
			return fmt.Errorf("data received")
		case DataChanError:
			if ContainsSuccessTerm(conf.SuccessTerm, SuccessTermError) {
				ctr.Logger.Info(ctx, "error received")
				return nil
			}
			ctr.Logger.Warn(ctx, "error received")
			return fmt.Errorf("error received")
		}
	case <-ctx.Done():
		ctr.Logger.Warn(ctx, "context cancelled")
		return fmt.Errorf("context cancelled")
	case <-timeout:
		if ContainsSuccessTerm(conf.SuccessTerm, SuccessTermTime) {
			ctr.Logger.Info(ctx, "timeout",
				logger.Value("time", breakTime))
			return nil
		}
		ctr.Logger.Warn(ctx, "timeout", logger.Value("time", breakTime))
		return fmt.Errorf("timeout exceeded")
	}

	return nil
}
