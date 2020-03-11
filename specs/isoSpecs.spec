# . is a special character and cannot appear between field names, spec names etc.
# field names within a spec are required to be unique
#
#
#format   := {fieldDef}={fieldSpecification}
#fieldDef := spec.{specName}.{messageName}.{fieldName}[.{positionInParent}.{childFieldName}]
#fieldSpefication (fixed size fields) := {fieldType}.{fieldEncoding}.{sizeSpec}[.constraints]
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
spec.TestSpec.Default Message.Message Type=fixed.ascii.size:4
spec.TestSpec.Default Message.Bitmap=bitmap.binary
spec.TestSpec.Default Message.Bitmap.2.Fixed2_ASCII=fixed.ascii.size:3.constraints{content:Numeric;}
spec.TestSpec.Default Message.Bitmap.3.Fixed3_EBCDIC=fixed.ebcdic.size:3.constraints{content:Alpha;}
spec.TestSpec.Default Message.Bitmap.4.Fixed4_BCD=fixed.bcd.size:3
spec.TestSpec.Default Message.Bitmap.5.Fixed5_BINARY=fixed.binary.size:3
#
## An example of a fixed field with embedded/nested subfields
spec.TestSpec.Default Message.Bitmap.6.FxdField6_WithSubFields=fixed.ascii.size:8
spec.TestSpec.Default Message.FxdField6_WithSubFields.1.SF6_1=fixed.ascii.size:4
spec.TestSpec.Default Message.SF6_1.1.SF6_1_1=fixed.ascii.size:2
spec.TestSpec.Default Message.SF6_1.1.SF6_1_2=fixed.ascii.size:2
spec.TestSpec.Default Message.FxdField6_WithSubFields.2.SF6_2=fixed.ascii.size:2
spec.TestSpec.Default Message.FxdField6_WithSubFields.3.SF6_3=fixed.ascii.size:2
#
## An example of a variable field with embedded/nested subfields
spec.TestSpec.Default Message.Bitmap.7.VarField7_WithSubFields=variable.ascii.binary.size:2
spec.TestSpec.Default Message.VarField7_WithSubFields.1.SF7_1=fixed.ascii.size:5
spec.TestSpec.Default Message.VarField7_WithSubFields.2.SF7_2=variable.bcd.ascii.size:2
spec.TestSpec.Default Message.VarField7_WithSubFields.3.SF7_3=fixed.ascii.size:5
#
# Rest of the fields
spec.TestSpec.Default Message.Bitmap.55.VarField55_BCD_BINARY=variable.bcd.binary.size:2
spec.TestSpec.Default Message.Bitmap.56.VarField56_BCD_ASCII=variable.bcd.ascii.size:2
spec.TestSpec.Default Message.Bitmap.57.VarField57_BINARY_EBCDIC =variable.binary.ebcdic.size:2
spec.TestSpec.Default Message.Bitmap.58.VarField58_EBCDIC_EBCDIC=variable.ebcdic.ebcdic.size:2
spec.TestSpec.Default Message.Bitmap.59.VarField59_EBCDIC_ASCII=variable.ebcdic.ascii.size:2
spec.TestSpec.Default Message.Bitmap.60.VarField60_EBCDIC_BINARY=variable.ebcdic.binary.size:3.constraints{minSize:8;maxSize:12;}
spec.TestSpec.Default Message.Bitmap.91.VarField91_ASCII_EBCDIC=variable.ascii.ebcdic.size:2.constraints{minSize:5;maxSize:15;content:Alpha;}
#
###
####
#Iso8583-MiniSpec
####
###
spec.Iso8583-MiniSpec.1100.Message Type=fixed.ascii.size:4
spec.Iso8583-MiniSpec.1100.Bitmap=bitmap.binary
spec.Iso8583-MiniSpec.1100.Bitmap.2.PAN=variable.ebcdic.ebcdic.size:2.constraints{content:Numeric;}
spec.Iso8583-MiniSpec.1100.Bitmap.3.Processing Code=fixed.ebcdic.size:6.constraints{content:Numeric;}
spec.Iso8583-MiniSpec.1100.Bitmap.4.Amount=fixed.ascii.size:12
spec.Iso8583-MiniSpec.1100.Bitmap.11.STAN=fixed.ascii.size:6
spec.Iso8583-MiniSpec.1100.Bitmap.38.Approval Code=fixed.ebcdic.size:6
spec.Iso8583-MiniSpec.1100.Bitmap.39.Action Code=fixed.ascii.size:3.constraints{content:Numeric;}
#
spec.Iso8583-MiniSpec.1420.Message Type=fixed.ascii.size:4
spec.Iso8583-MiniSpec.1420.Bitmap=bitmap.binary
spec.Iso8583-MiniSpec.1420.Bitmap.2.PAN=variable.ebcdic.ebcdic.size:2.constraints{content:Numeric;}
spec.Iso8583-MiniSpec.1420.Bitmap.3.Processing Code=fixed.ebcdic.size:6.constraints{content:Numeric;}
spec.Iso8583-MiniSpec.1420.Bitmap.4.Amount=fixed.ascii.size:12
spec.Iso8583-MiniSpec.1420.Bitmap.11.STAN=fixed.ascii.size:6
spec.Iso8583-MiniSpec.1420.Bitmap.38.Approval Code=fixed.ebcdic.size:6
spec.Iso8583-MiniSpec.1420.Bitmap.39.Action Code=fixed.ascii.size:3.constraints{content:Numeric;}
#