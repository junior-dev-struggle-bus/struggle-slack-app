exports.handler = function (event, context, callback) {

    const invokedByUser = event.body.user_name

    const res = {
       response_type: "in_channel",
        text: "http://devhumor.com/content/uploads/images/July2019/mvp_bugs.png",
        attachments: [
            {
                text: `Hey ${invokedByUser}! lowPriority is a simple Netlify function created by @Brett Hurst for the Struggle Buddy slack app. You should make it a highPriority to create your own! :D`
            }
        ]
    }

    callback(null, {
        statusCode: 200,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(res)
   })
}
