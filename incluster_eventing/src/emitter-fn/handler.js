const { v4: uuidv4 } = require('uuid');
const { SpanStatusCode } = require("@opentelemetry/api/build/src/trace/status");

module.exports = {
    main: async function (event, context) {
        let sanitisedData = sanitise(event.data)

        const msgId = uuidv4();
        const eventType = "sap.kyma.custom.acme.payload.sanitised.v1";
        const eventSource = "kyma";

        var eventOut=event.buildResponseCloudEvent(msgId,eventType,sanitisedData);
        eventOut.source=eventSource;
        eventOut.specversion="1.0";

        const span = event.tracer.startSpan('call-to-kyma-eventing');
        return await event.publishCloudEvent(eventOut)
            .then(resp => {
                if(resp.status!==204){
                    throw new Error("Unexpected response from eventing proxy");
                }
                span.addEvent("Event sent");
                span.setAttribute("event-type", eventType);
                span.setAttribute("event-source", eventSource);
                span.setAttribute("event-id", msgId);
                span.setStatus({code: SpanStatusCode.OK});
                return "Event sent";
            }).catch(err=> {
                console.error(err)
                span.setStatus({
                    code: SpanStatusCode.ERROR,
                    message: err.message,
                  });
                return err.message;
            }).finally(()=>{
                span.end();
            });
    }
}
let sanitise = (data)=>{
    console.log(`sanitising data...`)
    console.log(data)
    return data
}