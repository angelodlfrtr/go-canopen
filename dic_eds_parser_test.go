package canopen

import (
	"testing"
)

const TestEDSFile string = `
[Comments]

Lines=2

Line1=EDS file for AFE CANopen Slave

Line2=



[FileInfo]

FileName=E:\PROJEKT\CANOPEN_EDS\SEAFE.eds

FileVersion=1

FileRevision=0

EDSVersion=4.0

Description=EDS of the AFE

CreationTime=11:35AM

CreationDate=01-20-2004

CreatedBy=S.T.I.E.

ModificationTime=10:57AM

ModificationDate=04-02-2009

ModifiedBy=S.T.I.E.



[DeviceInfo]

Vendorname=Schneider Electric

VendorNumber=0x0200005A

ProductName=AFE_V1.0

ProductNumber=0x00414645

RevisionNumber=0x00010000

OrderCode=0

BaudRate_10=0

BaudRate_20=1

BaudRate_50=1

BaudRate_125=1

BaudRate_250=1

BaudRate_500=1

BaudRate_800=0

BaudRate_1000=1

SimpleBootUpMaster=0

SimpleBootUpSlave=1

Granularity=0

DynamicChannelsSupported=0

GroupMessaging=0

NrOfRXPDO=2

NrOfTXPDO=2

LSS_Supported=0

CompactPDO=0x00



[DummyUsage]

Dummy0001=1

Dummy0002=1

Dummy0003=1

Dummy0004=1

Dummy0005=1

Dummy0006=1

Dummy0007=1



[MandatoryObjects]

SupportedObjects=3

1=0x1000

2=0x1001

3=0x1018



[1000]

ParameterName=Device Type

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x00000000

PDOMapping=0

ObjFlags=0x0



[1001]

ParameterName=Error Register

ObjectType=0x7

DataType=0x0005

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x0

PDOMapping=0

ObjFlags=0x0



[1018]

ParameterName=Identity Object

SubNumber=4

ObjectType=0x8



[1018sub0]

ParameterName=Number of entries

ObjectType=0x7

DataType=0x0005

LowLimit=3

HighLimit=3

AccessType=ro

DefaultValue=3

PDOMapping=0

ObjFlags=0x0



[1018sub1]

ParameterName=Vendor ID

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x0200005A

PDOMapping=0

ObjFlags=0x0



[1018sub2]

ParameterName=Product code

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x00414645

PDOMapping=0

ObjFlags=0x0



[1018sub3]

ParameterName=Revision number

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x00010000

PDOMapping=0

ObjFlags=0x0



[OptionalObjects]

SupportedObjects=19

1=0x1003

2=0x1005

3=0x1008

4=0x100B

5=0x100C

6=0x100D

7=0x100E

8=0x100F

9=0x1014

10=0x1016

11=0x1017

12=0x1400

13=0x1401

14=0x1600

15=0x1601

16=0x1800

17=0x1801

18=0x1A00

19=0x1A01



[1003]

ParameterName=Pre-defined Error Field

SubNumber=2

ObjectType=0x8



[1003sub0]

ParameterName=Number of error

ObjectType=0x7

DataType=0x0005

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0

PDOMapping=0

ObjFlags=0x0



[1003sub1]

ParameterName=Standard Error Field

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0

PDOMapping=0

ObjFlags=0x0



[1005]

ParameterName=COB-ID SYNC message

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=rw

DefaultValue=0x80

PDOMapping=0

ObjFlags=0x0



[1008]

ParameterName=Device Name

ObjectType=0x7

DataType=0x0009

LowLimit=

HighLimit=

AccessType=const

DefaultValue=AFE

PDOMapping=0

ObjFlags=0x0



[100B]

ParameterName=NodeID

ObjectType=0x7

DataType=0x0005

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=

PDOMapping=0

ObjFlags=0x0



[100C]

ParameterName=Guard Time

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=rw

DefaultValue=0

PDOMapping=0

ObjFlags=0x0



[100D]

ParameterName=Life Time Factor

ObjectType=0x7

DataType=0x0005

LowLimit=

HighLimit=

AccessType=rw

DefaultValue=0

PDOMapping=0

ObjFlags=0x0



[100E]

ParameterName=Node Guarding Identifier

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=$NodeID + 0x700

PDOMapping=0

ObjFlags=0x0



[100F]

ParameterName=Number of SDO supported

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=1

PDOMapping=0

ObjFlags=0x0



[1014]

ParameterName=COB-ID Emergency message

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=$NodeID + 0x080

PDOMapping=0

ObjFlags=0x0



[1016]

ParameterName=Consumer Heartbeat Time

SubNumber=2

ObjectType=0x8



[1016sub0]

ParameterName=Number of entries

ObjectType=0x7

DataType=0x0005

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=1

PDOMapping=0

ObjFlags=0x0



[1016sub1]

ParameterName=Consumer Heartbeat Time

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=rw

DefaultValue=0

PDOMapping=0

ObjFlags=0x0



[1017]

ParameterName=Producer Heartbeat Time

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=rw

DefaultValue=0

PDOMapping=0

ObjFlags=0x0



[1400]

ParameterName=Receive PDO1 parameter

SubNumber=3

ObjectType=0x9



[1400sub0]

ParameterName=Number of entries

ObjectType=0x7

DataType=0x0005

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=2

PDOMapping=0

ObjFlags=0x0



[1400sub1]

ParameterName=COB-ID

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=rw

DefaultValue=$NodeID + 0x0200

PDOMapping=0

ObjFlags=0x0



[1400sub2]

ParameterName=Transmission type

ObjectType=0x7

DataType=0x0005

LowLimit=

HighLimit=

AccessType=rw

DefaultValue=255

PDOMapping=0

ObjFlags=0x0



[1401]

ParameterName=Receive PDO2 parameter

SubNumber=3

ObjectType=0x9



[1401sub0]

ParameterName=Number of entries

ObjectType=0x7

DataType=0x0005

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=2

PDOMapping=0

ObjFlags=0x0



[1401sub1]

ParameterName=COB-ID

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=rw

DefaultValue=$NodeID + 0x80000300

PDOMapping=0

ObjFlags=0x0



[1401sub2]

ParameterName=Transmission type

ObjectType=0x7

DataType=0x0005

LowLimit=

HighLimit=

AccessType=rw

DefaultValue=255

PDOMapping=0

ObjFlags=0x0



[1600]

ParameterName=Receive PDO1 mapping

SubNumber=5

ObjectType=0x8



[1600sub0]

ParameterName=Number of mapped objects

ObjectType=0x7

DataType=0x0005

LowLimit=0

HighLimit=4

AccessType=ro

DefaultValue=4

PDOMapping=0

ObjFlags=0x0



[1600sub1]

ParameterName=1.mapped object

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x30000110

PDOMapping=0

ObjFlags=0x3



[1600sub2]

ParameterName=2.mapped object

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x30000210

PDOMapping=0

ObjFlags=0x3



[1600sub3]

ParameterName=3.mapped object

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x30000310

PDOMapping=0

ObjFlags=0x3



[1600sub4]

ParameterName=4.mapped object

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x30000410

PDOMapping=0

ObjFlags=0x3



[1601]

ParameterName=Receive PDO2 mapping

SubNumber=5

ObjectType=0x8



[1601sub0]

ParameterName=Number of mapped objects

ObjectType=0x7

DataType=0x0005

LowLimit=0

HighLimit=4

AccessType=ro

DefaultValue=4

PDOMapping=0

ObjFlags=0x0



[1601sub1]

ParameterName=1.mapped object

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x30000510

PDOMapping=0

ObjFlags=0x3



[1601sub2]

ParameterName=2.mapped object

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x30000610

PDOMapping=0

ObjFlags=0x3



[1601sub3]

ParameterName=3.mapped object

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x30000710

PDOMapping=0

ObjFlags=0x3



[1601sub4]

ParameterName=4.mapped object

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x30000810

PDOMapping=0

ObjFlags=0x3



[1800]

ParameterName=Transmit PDO1 parameter

SubNumber=6

ObjectType=0x9



[1800sub0]

ParameterName=Number of entries

ObjectType=0x7

DataType=0x0005

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=5

PDOMapping=0

ObjFlags=0x0



[1800sub1]

ParameterName=COB-ID

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=rw

DefaultValue=$NodeID + 0x0180

PDOMapping=0

ObjFlags=0x0



[1800sub2]

ParameterName=Transmission type

ObjectType=0x7

DataType=0x0005

LowLimit=

HighLimit=

AccessType=rw

DefaultValue=255

PDOMapping=0

ObjFlags=0x0



[1800sub3]

ParameterName=Inhibit timer

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=rw

DefaultValue=300

PDOMapping=0

ObjFlags=0x0



[1800sub4]

ParameterName=reserved

ObjectType=0x7

DataType=0x0005

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=

PDOMapping=0

ObjFlags=0x0



[1800sub5]

ParameterName=Event Timer

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=rw

DefaultValue=1000

PDOMapping=0

ObjFlags=0x0



[1801]

ParameterName=Transmit PDO2 parameter

SubNumber=6

ObjectType=0x9



[1801sub0]

ParameterName=Number of entries

ObjectType=0x7

DataType=0x0005

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=5

PDOMapping=0

ObjFlags=0x0



[1801sub1]

ParameterName=COB-ID

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=rw

DefaultValue=$NodeID + 0x80000280

PDOMapping=0

ObjFlags=0x0



[1801sub2]

ParameterName=Transmission type

ObjectType=0x7

DataType=0x0005

LowLimit=

HighLimit=

AccessType=rw

DefaultValue=255

PDOMapping=0

ObjFlags=0x0



[1801sub3]

ParameterName=Inhibit timer

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=rw

DefaultValue=300

PDOMapping=0

ObjFlags=0x0



[1801sub4]

ParameterName=reserved

ObjectType=0x7

DataType=0x0005

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=

PDOMapping=0

ObjFlags=0x0



[1801sub5]

ParameterName=Event Timer

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=rw

DefaultValue=1000

PDOMapping=0

ObjFlags=0x0



[1A00]

ParameterName=Transmit PDO1 mapping

SubNumber=5

ObjectType=0x8



[1A00sub0]

ParameterName=Number of mapped objects

ObjectType=0x7

DataType=0x0005

LowLimit=0

HighLimit=4

AccessType=ro

DefaultValue=4

PDOMapping=0

ObjFlags=0x0



[1A00sub1]

ParameterName=1.mapped object

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x30100110

PDOMapping=0

ObjFlags=0x3



[1A00sub2]

ParameterName=2.mapped object

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x30100210

PDOMapping=0

ObjFlags=0x3



[1A00sub3]

ParameterName=3.mapped object

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x30100310

PDOMapping=0

ObjFlags=0x3



[1A00sub4]

ParameterName=4.mapped object

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x30100410

PDOMapping=0

ObjFlags=0x3



[1A01]

ParameterName=Transmit PDO2 mapping

SubNumber=5

ObjectType=0x8



[1A01sub0]

ParameterName=Number of mapped objects

ObjectType=0x7

DataType=0x0005

LowLimit=0

HighLimit=4

AccessType=ro

DefaultValue=4

PDOMapping=0

ObjFlags=0x0



[1A01sub1]

ParameterName=1.mapped object

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x30100510

PDOMapping=0

ObjFlags=0x3



[1A01sub2]

ParameterName=2.mapped object

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x30100610

PDOMapping=0

ObjFlags=0x3



[1A01sub3]

ParameterName=3.mapped object

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x30100710

PDOMapping=0

ObjFlags=0x3



[1A01sub4]

ParameterName=4.mapped object

ObjectType=0x7

DataType=0x0007

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=0x30100810

PDOMapping=0

ObjFlags=0x3



[ManufacturerObjects]

SupportedObjects=2

1=0x3000

2=0x3010



[3000]

ParameterName=controlword and target values

SubNumber=9

ObjectType=0x8



[3000sub0]

ParameterName=number of target values

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=8

PDOMapping=0

ObjFlags=0x0



[3000sub1]

ParameterName=controlword

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=rww

DefaultValue=

PDOMapping=1

ObjFlags=0x0



[3000sub2]

ParameterName=ref_value 1

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=rww

DefaultValue=

PDOMapping=1

ObjFlags=0x0



[3000sub3]

ParameterName=ref_value 2

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=rww

DefaultValue=

PDOMapping=1

ObjFlags=0x0



[3000sub4]

ParameterName=ref_value 3

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=rww

DefaultValue=

PDOMapping=1

ObjFlags=0x0



[3000sub5]

ParameterName=ref_value 4

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=rww

DefaultValue=

PDOMapping=1

ObjFlags=0x0



[3000sub6]

ParameterName=ref_value 5

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=rww

DefaultValue=

PDOMapping=1

ObjFlags=0x0



[3000sub7]

ParameterName=ref_value 6

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=rww

DefaultValue=

PDOMapping=1

ObjFlags=0x0



[3000sub8]

ParameterName=ref_value 7

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=rww

DefaultValue=

PDOMapping=1

ObjFlags=0x0



[3010]

ParameterName=statusword and actual values

SubNumber=9

ObjectType=0x8



[3010sub0]

ParameterName=number of actual values

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=8

PDOMapping=0

ObjFlags=0x0



[3010sub1]

ParameterName=statusword

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=

PDOMapping=1

ObjFlags=0x0



[3010sub2]

ParameterName=act_value 1

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=

PDOMapping=1

ObjFlags=0x0



[3010sub3]

ParameterName=act_value 2

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=

PDOMapping=1

ObjFlags=0x0



[3010sub4]

ParameterName=act_value 3

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=

PDOMapping=1

ObjFlags=0x0



[3010sub5]

ParameterName=act_value 4

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=

PDOMapping=1

ObjFlags=0x0



[3010sub6]

ParameterName=act_value 5

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=

PDOMapping=1

ObjFlags=0x0



[3010sub7]

ParameterName=act_value 6

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=

PDOMapping=1

ObjFlags=0x0



[3010sub8]

ParameterName=act_value 7

ObjectType=0x7

DataType=0x0006

LowLimit=

HighLimit=

AccessType=ro

DefaultValue=

PDOMapping=1

ObjFlags=0x0
`

func TestDicEDSParse(t *testing.T) {
	// Parse file
	dic, err := DicEDSParse([]byte(TestEDSFile))

	if err != nil {
		t.Fatal(err)
	}

	t.Log(dic)
}
