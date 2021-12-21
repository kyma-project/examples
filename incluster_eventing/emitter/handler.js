module.exports = {
    main: function (event, context) {
        console.log(event.data)
        return 'Hello Serverless'
    }
}