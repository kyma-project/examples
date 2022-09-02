const { v4: uuidv4 } = require('uuid');
const { SpanStatusCode } = require("@opentelemetry/api/build/src/trace/status");

module.exports = {
    main: async function (event, context) {

        const msgId = uuidv4();
        const eventType = process.env['eventtype']
        const eventSource = process.env['eventsource']
        const eventSpecVersion = process.env['eventspecversion']

        var eventOut=event.buildResponseCloudEvent(msgId,eventType,event.data);
        eventOut.source=eventSource;
        eventOut.specversion=eventSpecVersion

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
                return "Event sent : "+JSON.stringify(event.data);
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
