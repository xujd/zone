package app

import "github.com/spf13/pflag"

// ServerOption define option of server in flags
type ServerOption struct {
	AddrPort   string
	DbHost     string
	DbPort     int
	DbUser     string
	DbPassword string
	DbName     string
	FileDir    string
}

// NewServerOption create a ServerOption object
func NewServerOption() *ServerOption {
	s := ServerOption{
		AddrPort:   ":1323",
		DbHost:     "127.0.0.1",
		DbPort:     5432,
		DbUser:     "postgres",
		DbPassword: "12345678",
		DbName:     "cmkit",
		FileDir:    "./webfiles",
	}

	return &s
}

// AddFlags add flags
func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.AddrPort, "addrport", ":1323", "The ip address and port for the serve on")
	fs.StringVar(&s.DbHost, "dbhost", "127.0.0.1", "The db address")
	fs.IntVar(&s.DbPort, "dbport", 5432, "The db port")
	fs.StringVar(&s.DbUser, "dbuser", "postgres", "The db user")
	fs.StringVar(&s.DbPassword, "dbpassword", "12345678", "The db password")
	fs.StringVar(&s.DbName, "dbname", "cmkit", "The db name")
	fs.StringVar(&s.FileDir, "filedir", "./webfiles", "The filedir for upload")
}
