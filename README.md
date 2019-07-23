# Struggle Buddy, a JDSB Slack App

## Contributing
1. [Fork this repo](https://help.github.com/en/articles/fork-a-repo#fork-an-example-repository).
2. Add your beautiful JS code to the [functions/ folder](https://github.com/junior-dev-struggle-bus/struggle-slack-app/tree/staging/functions).
3. Register your function with the [public/slack-app-registry.json](https://github.com/junior-dev-struggle-bus/struggle-slack-app/blob/staging/public/slack-app-registry.json) under the "functions" entry. Here's a [real-life reference](https://github.com/junior-dev-struggle-bus/struggle-slack-app/commit/0ff6622028e87e0729fc97feda3c4080f4606eb1). Otherwise, the general form is below:
```json
   "yourFunctionName" : {
        "usage" : "/struggle yourFunctionName faveArg lessFaveArg worstArg -OR- anything, this gets displayed to the user.",
	"description" : "The best description you can muster.",
	"manual" : "A link to site detailing more on what your function does."
   }
```
   The above will tell the backend to pick-up this function for deployment automatically. You don't worry about anything else.
   
4. Git add, commit, and push your files to your forked repo. [Here's a tutorial if you're not familiar](https://www.atlassian.com/git).
5. [Create a Pull Request](https://help.github.com/en/articles/creating-a-pull-request) against this repo.
6. Ping a repo owner or collaborator to review your code either here or on our [JDSB Slack workspace in the #struggle-app-project channel](https://join.slack.com/t/jdsb/shared_invite/enQtNzA0NTY3OTE2ODg3LTE5ZTE5ODI5YmE5YTUzN2UyOWUxZmM1ZDZlNDliZTgxYTg0ODRlMmM3OThkY2JlZDRlNjIzYmJiMjNjNDBjOWQ).
7. Once approved and merged, yourFunctionName will be testable in the [JDSB Testing Ground Slack workspace](https://join.slack.com/t/jdsb-wrecking-ball/shared_invite/enQtNjgyMjA3NzU4MzIyLWMyNjIyZDY3ZDkwMTdiY2VlNDhlNDg2YTYyODQ3ZjRlZjA1NTZiNmNhZjcyNDM5MDhiNDliMmFmMzExOTJiNTk).
8. After some QA testing, we'll merge your code into the production branch where it'll be available for all to use in the [JDSB Slack workspace](https://join.slack.com/t/jdsb\/shared_invite/enQtNzA0NTY3OTE2ODg3LTE5ZTE5ODI5YmE5YTUzN2UyOWUxZmM1ZDZlNDliZTgxYTg0ODRlMmM3OThkY2JlZDRlNjIzYmJiMjNjNDBjOWQ).

## Appendix

### Deployment Infrastructure Design
[Design Document](https://docs.google.com/document/d/16dnOpjcIh5SccWKhtK9-ZRwT3qARPdJVvA3x1K9aIrY/edit?usp=sharing)
