#Introduction

This software is intended to receive all measurements from an OpenAir
device for testing and calibration. This is more data than would be
useful in most cases, hence the name 'firehose'.


## Running

This software is a commandline utility. Precompiled versions of the
binary for Linux, MacOS and Windows plattform are available on the
github releases page.

The software is a command line program with the following options

````
$ ./firehose -help

Usage of ./firehose:

  -a string

  -addr string

        address for server to listen on (default ":7531")

  -h    

  -help

        print usage

  -o string

  -outfile string

        filename to save output to (default "-")

````

By default the software listens to port 7531 on all interfaces, using
this option, it can be restricted to a single IP address, e.g. the
address assigned by OpenAir if connecting to the AP provided by the
device.

The software receives sensor readings (currently as a TCP Stream of JSON
packets) and outputs a CSV file for use in Excel, etc. By default the
CSV entries are printed to STDOUT, but a filename to save the csv can be
provided using the outfile parameter.

The program can receive and process output from multiple OpenAir
devices. These devices must be configured to use the same server using
the firehose_addr configuration parameter.

The format of the CSV is as follows:

````
server_time,timestamp,device_id,tag,value(hex),value(decimal),tag_annotation,value_annotation
1559312722,0,esp32_0AAEAC,0xffffffff,0x00000000,0,OA_Network_Events,CONNECT                                                                                                                                        
1559312722,992,esp32_0AAEAC,0x00000009,0x0004d243,315971,OA_BME_Pressure_Raw,(raw 315971)                                                                                                                          
1559312722,992,esp32_0AAEAC,0x0000000a,0x00018e28,101928,OA_BME_Pressure,1019.28 hPa                                                                                                                               
1559312722,992,esp32_0AAEAC,0x0000000b,0x000801dc,524764,OA_BME_Temp_Raw,(raw 524764)                                                                                                                              
1559312722,992,esp32_0AAEAC,0x0000000c,0x00048d3b,298299,OA_BME_Temp,25.15 C
...
````


Field | Description
------|------------
server_time | Time the package was received by the server (unix timestamp, seconds since 1970-01-01)
timestamp | OpenAir timestamp of sensor reading on the device, seconds since device boot
device_id | String identifying the OpenAir device
tag | Tag identifying the source of the reading. For convenience this value is annotated in the field tag_annotation. A list of all tags is provided in a table below.
value(hex) | All sensor reading are currently stored as unsigned 32bit values. This field is a hex representation of the value.
value(dec) | This field is a decimal representation of the sensor value. Depending on the nature of the data returned by the sensor, this may not be a sensible representation.
tag_annotation | Human readable description of the tag
value_annotation | Human readable interpretation of the value for convenience.
OA_Network_Events | Value CONNECT and DISCONNECT signify time the OpenAir device established a connection or the connection was lost.
OA_AlphaCalc_1 … OA_Alpha_Calc_4 | Precalculated (server side) voltage of the ADC reading. This value is calculated by the formula:

````
V = (ALPHA(X-1) - ALPHA(X)) * ADC_CONST

For X in 2, 4, 6, 8
ADC_CONST = 0.000031356811523
```1`
