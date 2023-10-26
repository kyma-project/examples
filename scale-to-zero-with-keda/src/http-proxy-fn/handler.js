const { SpanStatusCode } = require("@opentelemetry/api/build/src/trace/status");

module.exports = {
    main: async function (event, context) {

        const eventType = process.env['eventtype']
        const eventSource = process.env['eventsource']

        const span = event.tracer.startSpan('call-to-kyma-eventing');

        return await event.emitCloudEvent(eventType, eventSource, event.data)
        .then(resp => {
            console.log(resp.status);
            span.addEvent("Event sent");
            span.setAttribute("event-type", eventType);
            span.setAttribute("event-source", eventSource);
            span.setStatus({code: SpanStatusCode.OK});
            return "Event sent: "+JSON.stringify(event.data);
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
