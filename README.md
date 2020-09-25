# squaregeorge
Sqauregeorge identifies the responsible mailservers behind massive amounts of mail addresses.

squaregeorge takes a list of mail addresses as argument, tries to parse them, and creates a 
list of all mail servers that are involved in handling mails sent to all of these addresses. 
This information is obtained by querying all found domains for MX records.

The intended use case is to identify all servers that need to be reconfigured in order to
allow desirable mass-mailings to reach the destined recipients. 
