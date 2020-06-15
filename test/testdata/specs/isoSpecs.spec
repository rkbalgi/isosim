# . is a special character and cannot appear between field names, spec names etc.
# field names within a spec are required to be unique
# Message Type (or can also be called MTI) and Bitmap are special fields and their names should'nt be changed
#
#
###
#specDef := spec.{specName}.{Id}=UniqueSpecId
#msgDef :={UniqueSpecId|SpecName}.{MessageName}.Id={UniqueMsgIdInSpec}
###
#fieldFormat := {fieldDef}={fieldSpecification}
#fieldDef := {specId|specName}.{msgId|msgName}.{fieldName}[.{positionInParent}.{childFieldName}]
#fieldSpefication (fixed size fields) := {fieldId}.{fieldType}.{fieldEncoding}.{sizeSpec}[.constraints]
#sizeSpec := size:[0-9]+
#fieldSpefication (variable size fields) := {fieldType}.{lengthEncoding}.{dataEncoding}.{lengthEncodingSizeSpec}[.constraints}
#lengthEncodingSizeSpec := size:[0-9]+
#constraints:=constraints'{' [content:{Numeric|AlphaNumeric}];[minSize:[0-9]+];[maxSize:[0-9]+]'}'
##
#
#
###
####
#TestSpec
###
####
spec.TestSpec.Id=1
1.Default Message.Id=1
####
1.1.Message Type.1=fixed.ascii.size:4
1.1.Bitmap.2=bitmap.binary
1.1.Bitmap.2.Fixed2_ASCII.3=fixed.ascii.size:3.constraints{content:Numeric;}
1.1.2.3.Fixed3_EBCDIC.4=fixed.ebcdic.size:3.constraints{content:Alpha;}
1.1.Bitmap.4.Fixed4_BCD.5=fixed.bcd.size:3
1.1.Bitmap.5.Fixed5_BINARY.6=fixed.binary.size:3
#
## An example of a fixed field with embedded/nested subfields
1.1.2.6.FxdField6_WithSubFields.7=fixed.ascii.size:8
1.1.7.1.SF6_1.8=fixed.ascii.size:4
1.1.8.1.SF6_1_1.9=fixed.ascii.size:2
1.1.SF6_1.1.SF6_1_2.10=fixed.ascii.size:2
1.1.FxdField6_WithSubFields.2.SF6_2.11=fixed.ascii.size:2
1.1.FxdField6_WithSubFields.3.SF6_3.12=fixed.ascii.size:2
#
## An example of a variable field with embedded/nested subfields
1.1.Bitmap.7.VarField7_WithSubFields.13=variable.ascii.binary.size:2
1.1.VarField7_WithSubFields.1.SF7_1.14=fixed.ascii.size:5
1.1.VarField7_WithSubFields.2.SF7_2.15=variable.bcd.ascii.size:2
1.1.VarField7_WithSubFields.3.SF7_3.16=fixed.ascii.size:5
#
# Rest of the fields
1.1.Bitmap.55.VarField55_BCD_BINARY.17=variable.bcd.binary.size:2
1.1.Bitmap.56.VarField56_BCD_ASCII.18=variable.bcd.ascii.size:2
1.1.Bitmap.57.VarField57_BINARY_EBCDIC.19=variable.binary.ebcdic.size:2
1.1.Bitmap.58.VarField58_EBCDIC_EBCDIC.20=variable.ebcdic.ebcdic.size:2
1.1.Bitmap.59.VarField59_EBCDIC_ASCII.21=variable.ebcdic.ascii.size:2
1.1.Bitmap.60.VarField60_EBCDIC_BINARY.22=variable.ebcdic.binary.size:3.constraints{minSize:8;maxSize:12;}
1.1.Bitmap.91.VarField91_ASCII_EBCDIC.23=variable.ascii.ebcdic.size:2.constraints{minSize:5;maxSize:15;content:Alpha;}
#
###
####
#Iso8583-MiniSpec
####
###
spec.Iso8583-MiniSpec.Id=2
2.1100.Id=1
Iso8583-MiniSpec.1420.Id=2
###
2.1.Message Type.1=fixed.ascii.size:4
2.1.Bitmap.2=bitmap.binary
2.1.Bitmap.2.PAN.3=variable.ebcdic.ebcdic.size:2.constraints{content:Numeric;}
2.1.Bitmap.3.Processing Code.4=fixed.ebcdic.size:6.constraints{content:Numeric;}
2.1.2.4.Amount.5=fixed.ascii.size:12
2.1.2.11.STAN.6=fixed.ascii.size:6
2.1.2.38.Approval Code.7=fixed.ebcdic.size:6
2.1.Bitmap.39.Action Code.8=fixed.ascii.size:3.constraints{content:Numeric;}
#
2.2.Message Type.1=fixed.ascii.size:4
2.2.Bitmap.2=bitmap.binary
2.2.Bitmap.2.PAN.3=variable.ebcdic.ebcdic.size:2.constraints{content:Numeric;}
2.2.Bitmap.3.Processing Code.4=fixed.ebcdic.size:6.constraints{content:Numeric;}
2.2.Bitmap.4.Amount.5=fixed.ascii.size:12
2.2.2.11.STAN.6=fixed.ascii.size:6
2.2.Bitmap.37.Retrieval Reference Number.7=variable.ascii.ascii.size:2
2.2.2.38.Approval Code.8=fixed.ebcdic.size:6
2.2.Bitmap.39.Action Code.9=fixed.ascii.size:3.constraints{content:Numeric;}
#