# . is a special character and cannot appear between field names, spec names etc.
#
#
#format   := {fieldDef}={fieldSpecification}
#fieldDef := spec.{specName}.{messageName}.{fieldName}[.{positionInParent}.{childFieldName}]
#fieldSpefication (fixed size fields) := {fieldType}.{fieldEncoding}.{sizeSpec}[.constraints]
#sizeSpec := size:[0-9]+
#fieldSpefication (variable size fields) := {fieldType}.{lengthEncoding}.{dataEncoding}.{lengthEncodingSizeSpec}[.constraints}
#lengthEncodingSizeSpec := size:[0-9]+
#constraints:=constraints'{' [content:{Numeric|AlphaNumeric}];[minSize:[0-9]+];[maxSize:[0-9]+]'}'
#
spec.TestSpec.Default Message.Message Type=fixed.ascii.size:4
spec.TestSpec.Default Message.Bitmap=bitmap.binary
spec.TestSpec.Default Message.Bitmap.2.Fixed ASCII=fixed.ascii.size:3.constraints{content:Numeric;}
spec.TestSpec.Default Message.Bitmap.3.Fixed EBCDIC=fixed.ebcdic.size:3.constraints{content:Alpha;}
spec.TestSpec.Default Message.Bitmap.4.Fixed BCD=fixed.bcd.size:3
spec.TestSpec.Default Message.Bitmap.5.Fixed BINARY=fixed.binary.size:3
spec.TestSpec.Default Message.Bitmap.55.Var BCD/BINARY=variable.bcd.binary.size:2
spec.TestSpec.Default Message.Bitmap.56.Var BCD/ASCII=variable.bcd.ascii.size:2
spec.TestSpec.Default Message.Bitmap.57.Var BINARY/EBCDIC =variable.binary.ebcdic.size:2
spec.TestSpec.Default Message.Bitmap.58.Var EBCDIC/EBCDIC=variable.ebcdic.ebcdic.size:2
spec.TestSpec.Default Message.Bitmap.59.Var EBCDIC/ASCII=variable.ebcdic.ascii.size:2
spec.TestSpec.Default Message.Bitmap.60.Var EBCDIC/BINARY=variable.ebcdic.binary.size:3.constraints{minSize:8;maxSize:12;}
spec.TestSpec.Default Message.Bitmap.91.Var ASCII/EBCDIC=variable.ascii.ebcdic.size:2.constraints{minSize:5;maxSize:15;content:Alpha;}
#
####
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