package packet

var ExportReadPublishFixedHeader = (*MQTTReader).readPublishFixedHeader
var ExportReadPublishVariableHeader = (*MQTTReader).readPublishVariableHeader

var ExportReadVariableConnectHeader = (*MQTTReader).readConnectVariableHeader
var ExportReadConnectPayload = (*MQTTReader).readConnectPayload
