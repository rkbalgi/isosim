package websim

const (
	URLAllSpecs      = "/iso/v1/specs"
	URLMessages4Spec = "/iso/v1/msgs/" //{specId}

	URLGetMessageTemplate = "/iso/v1/template/" //{specId}/{msgId}
	URLParseTrace         = "/iso/v1/parse/"    //{specId}/{msgId}/
	URLParseTraceExternal = "/iso/v1/parse/external"
	URLSendMessageToHost  = "/iso/v1/send"
	URLSaveMsg            = "/iso/v1/save"
	URLLoadMsg            = "/iso/v1/loadmsg"
)
