module.exports = {
    main: async function (event, context) {
        return await processData(event.data);
    }
}
let processData = (data)=>{
    console.log(`Processing ...`);
    console.log(data);          
    return new Promise((resolve, reject) => {
        setTimeout(() => {
          resolve(`Done processing`);
        }, 1000)
    });

}