//Package storage is the package that handles all the data persistency.
//It's responsable to encrypt and decrypt the data before it reaches the disk.
//All the abstractions are developed to easy the exchange of pieces and testing
//
//Every persistency new "driver" should implement the interface interfaces.Job.
//
//To select the proper driver all you need to do is provide a crypt.Cipher
//and an connection URL.
//
//Valid conneciton URL's are:
//
// memory://
// mysql://user:password@host/database?parseTime=true
//
package storage
