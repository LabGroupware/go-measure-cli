package cmdreq

type CmdReq struct {
	CmdExecute func() (string, error)
}
