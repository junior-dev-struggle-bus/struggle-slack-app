exports.handler = function (event, context, callback) {

    const res = {
        response_type: "ephemeral",
        text: "http://devhumor.com/content/uploads/images/July2019/mvp_bugs.png"
    }

    callback(null, {
        statusCode: 200,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(res)
   })
}
