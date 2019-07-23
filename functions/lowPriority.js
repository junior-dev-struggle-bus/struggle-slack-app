exports.handler = function (event, context, callback) {

    const res = {
        response_type: "in_channel",
        text: "http://devhumor.com/content/uploads/images/July2019/mvp_bugs.png"
    }

    callback(null, {
        statusCode: 200,
        headers: { "Content-Type": "application/json" },
        body: launchesNextWeek
    })
}