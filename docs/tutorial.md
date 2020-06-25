## <b>Tutorial</b>

> Welcome! This is the tutorial for creating your own connetor using this framework. Perform the following changes in the framework code base to create your very own connector:

#### 1) Source Configuration File

* Change the name of the source config file to the source name and add the required configurations and credentials to it or you can simply create your own *source* configuration file.

#### 2) Source.go

Change the name of the file to *source name* and make the following amendments:

* Change the source configuration structure to store the specified source configurations.
* Add print statements in the Load<Source>Property function to print the specified configurations.
* Add the code to connect to source and create an instance, create/fetch back-up data or files, and create a reader to the backup file/data.

#### 3) Store.go

Following changes need to be made in the store.go file:

* Change the corresponding variable and flag names.
* Process the upload file name to convert to a standard and less complex form, if required.
* Made changes in the code fragment calling the upload function as per the arguments you wish to pass. Only the *reader* and *file path/name* arguments should be changed.

#### 4) Storj.go

The following changes need to be made only in the upload function:

* Change the arguments in the function definition as per the arguments passed from *store.go*.
* If reader is not passed as an argument to call the upload function, add the code fragment to create one. Remember to close the reader after the upload is committed.
* In case you want to implement section uploading, the corresponding code snippet has been provided in a commented block for the same. Uncomment the same and use it for the purpose.
* For uplaoding a byte array(buffer), a commented block has been provided for that also. Uncomment the same and use it for the purpose.

#### 5) Change the connector name in *root.go* and *main.go* files.

#### 6) Create a *go.mod* file for the respective connector.
