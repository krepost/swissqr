# QR Bills in Switzerland

Package swissqr creates a QR invoice form as described in version 2.0
of the document “[Schweizer Implementation Guidelines QR-Rechnung]
(https://www.paymentstandards.ch/dam/downloads/ig-qr-bill-de.pdf)”,
dated 15 November 2018, together with version 1.2 of the document
“[Syntaxdefinition der Rechnungsinformationen (S1) bei der QR-Rechnung]
(https://www.swiss-qr-invoice.org/downloads/qr-bill-s1-syntax-de.pdf)”,
dated 23 November 2018. This is not an officially supported Google product.

For example usage, please consult the `example_*_test.go` files. The general
pattern is: first initialize a `struct Payload` with the invoice content. Next,
validate that the payload is correct by calling the `Validate()` method on the
payload. Last but not least, create the actual invoice and store it in a PDF
document. When serializing the payload, it is a precondition that the payload
be valid.

The implementation is a best-effort to satisfy the standard to the letter as
well as the intention of the standard. Some points that were not clear from the
standard have been clarified based on the validation tool available at
[https://www.swiss-qr-invoice.org/]. Most importantly, there is an additional
check to verify that only the characters specified in the _Swiss Implementation
Guidelines for Customer-Bank Messages Credit Transfer_ are used in the invoice:

```
abcdefghijklmnopqrstuvwxyz
ABCDEFGHIJKLMNOPQRSTUVWXYZ
0123456789.,:'+-/()?
 !\"#%&*;<>÷=@_$£[]{}\`´~
àáâäçèéêëìíîïñòóôöùúûüýß
ÀÁÂÄÇÈÉÊËÌÍÎÏÒÓÔÖÙÚÛÜÑ
```
