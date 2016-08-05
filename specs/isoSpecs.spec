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
spec.ISO8583.1100.Message Type=fixed.ascii.size:4
spec.ISO8583.1100.Bitmap=bitmap.binary
spec.ISO8583.1100.Bitmap.2.Fixed ASCII=fixed.ascii.size:3.constraints{content:Numeric;}
spec.ISO8583.1100.Bitmap.3.Fixed EBCDIC=fixed.ebcdic.size:3.constraints{content:Alpha;}
spec.ISO8583.1100.Bitmap.4.Fixed BCD=fixed.bcd.size:3
spec.ISO8583.1100.Bitmap.5.Fixed BINARY=fixed.binary.size:3
spec.ISO8583.1100.Bitmap.55.Var BCD/BINARY=variable.bcd.binary.size:2
spec.ISO8583.1100.Bitmap.56.Var BCD/ASCII=variable.bcd.ascii.size:2
spec.ISO8583.1100.Bitmap.57.Var BINARY/EBCDIC =variable.binary.ebcdic.size:2
spec.ISO8583.1100.Bitmap.58.Var EBCDIC/EBCDIC=variable.ebcdic.ebcdic.size:2
spec.ISO8583.1100.Bitmap.59.Var EBCDIC/ASCII=variable.ebcdic.ascii.size:2
spec.ISO8583.1100.Bitmap.60.Var EBCDIC/BINARY=variable.ebcdic.binary.size:3.constraints{minSize:8;maxSize:12;}
spec.ISO8583.1100.Bitmap.91.Var ASCII/EBCDIC=variable.ascii.ebcdic.size:2.constraints{minSize:5;maxSize:15;content:Alpha;}
#
####