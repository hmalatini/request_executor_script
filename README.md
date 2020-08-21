# Request Executor Script
This script help when you have to make a REST call to the same endpoint, but with a lot of different data. So, you can parametrized the request, and let the golang workers to do all the job.
![](https://miro.medium.com/proxy/1*MKoHL3nNl3ZxDVMNBEZFTw.png)
## Execution
```bash
go run main.go
```
### Optional arguments
- **--help**: Help of the script
- **--config**: Path to the config file. For default, it will take the `config.yml` file inside the config folder.
## Configuration File
The configuration file could be divided in the following sections
### Data
For now, it is only allowed to provide **CSV** input files.
- **file**: Path to input file. Is where the script will collect the data, for parametrize the request.
- **header**: true/false. If the csv file contains header, put this flag in true, else, put in false
### Processor
For now, it can only process http requests.
Here is where the request is parametrized and made. So, here there are the posible values:
- **method**: REST method to execute.
- **baseUrl**: Base URL or hostname of the request
- **headers**: Headers of the request
- **uri**: Uri of the request. Use %N for parametrize a value. Where N is the position of the column starting from 1 to N
- **body**: Body of the request. Use %N for parametrize a value. Where N is the position of the column starting from 1 to N
### Result
Here is where the result of the execution is stored. So, it will be stored in a file. For now, it can only store in a CSV file.
The possible values for this config are:
- **file**: Path to result (output) file. If the file specified not exists, it will create a new one. If nothing is provided here, it will create a file in the current dir with the following name: `result-${timestamp}.csv`
- **append_results**: If the file specified above already exist and this flag is true, it will append the results to that file. If the file exists and this flag is true, it will delete the file and create a new one.
- **write_success**: Save on result file the requests that was successful
- **write_fails**: Save on result file the request that have an error

_Tip: If you set `write_success` on `false` and `write_fails` on `true`, you can re-run the script with only the failed data, specifying as input file, the result file of the previous execution._
### Executor
The executor is in charge to send to execute the request. You can configure the following parameters:
- **type**: Type of executor. Posible values are:
    - parallel: This is for execute request with goroutines
    - example: This executor will not execute any request, just print how the request could be parsed, for debug purposes
- **routines**: Number of go routines that will be executing
- **flush_number**: Chunk numbers
### Logger
The logger can show what are happening in the script while it is executing.
- **level**: TRACE, DEBUG, INFO, WARNING, ERROR, OFF.. If nothing is specified, the level `INFO` will be applied

---
![bitmoji](https://sdk.bitmoji.com/render/panel/0e5fd403-52c2-4e30-842e-19f331349c0b-e1b374fb-68ff-4b8b-bb46-a837a24fb984-v1.png?transparent=1&palette=1&widht=246)

For any suggestion, you can make a fork and then create a PR, or you can directly propose it at my [email](mailto:nani93@gmail.com).
