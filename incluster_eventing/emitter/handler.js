const { v4: uuidv4 } = require('uuid');
module.exports = {
    main: function (event, context) {
        let sanitisedData = sanitise(event.data)
        var eventOut=event.buildResponseCloudEvent(uuidv4(),"sap.kyma.custom.acme.payload.sanitised.v1",sanitisedData);
        eventOut.source="kyma"
        eventOut.specversion="1.0"
        event.publishCloudEvent(eventOut);
        console.log(`Payload pushed to sap.kyma.custom.acme.payload.sanitised.v1`,eventOut)
        return eventOut;
    }
}
let sanitise = (data)=>{
    console.log(`sanitising data...`)
    console.log(data)
    return data
}