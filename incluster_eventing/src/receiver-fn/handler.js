module.exports = {
    main: function (event, context) {
        store(event.data)
        return 'Stored'
    }
}
let store = (data)=>{
    console.log(`storing data...`)
    console.log(data)
    return data
}