exports.handler = function (event, context, callback) {

    const launchesNextWeek = `<html><body><img alt="Funny laughs and good times" src="http://devhumor.com/content/uploads/images/July2019/mvp_bugs.png"/><</body></html>`

    callback(null, {
        statusCode: 200,
        body: launchesNextWeek
    })
}