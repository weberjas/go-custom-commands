package cec_messages_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "splunk.com/go-custom-commands/cec_messages"
)

var _ = Describe("CEC Messages", func() {

	validActionMessage := "chunked 1.0,42,150\n" +
		"{ \"action\": \"execute\", \"finished\": false }\n" +
		"_raw,_time,index,sourcetype,source,host\n" +
		"\"you suck\",1234567890,main,silly,udp,mrt.splunk.com\n" +
		"\"you still suck\",1234567891,main,silly,udp,mrt.splunk.com"

	validGetinfoMessage := "chunked 1.0,23,0\n" +
		"{ \"action\": \"getinfo\" }"

	validComplexGetinfoMessage := "chunked 1.0,444,0\n" +
		"{ \"action\": \"getinfo\"," +
		"\"preview\": false," +
		"\"streaming_command_will_restart\": false," +
		"\"searchinfo\": {" +
		"\"args\": [ \"arg1\", \"arg2\" ]," +
		"\"raw_args\": [ \"rawarg1\", \"rawarg2\" ]," +
		"\"dispatch_dir\": \"/tmp\"," +
		"\"sid\": \"unittest\"," +
		"\"app\": \"app\"," +
		"\"username\": \"weberjas\"," +
		"\"owner\": \"jasonw\"," +
		"\"session_key\": \"123ABC\"," +
		"\"splunkd_uri\": \"/custom\"," +
		"\"splunk_version\": \"8.8.8\"," +
		"\"search\": \"SPL String\"," +
		"\"command\": \"commandName\"," +
		"\"maxresultrows\": 100," +
		"\"earliest_time\": 123456," +
		"\"latest_time\": 987654" +
		"}" +
		"}"

	It("correctly parses CEC protocol 'action' message from splunk", func() {

		cecMessage := SplunkCecMessage{}
		parsedMessage, err := cecMessage.Parse(validActionMessage)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(parsedMessage.ProtocolVersion).Should(Equal("chunked 1.0"))
		Expect(parsedMessage.MetadataLength).Should(Equal(42))
		Expect(parsedMessage.DataLength).Should(Equal(150))
		Expect(parsedMessage.MetaData.Action).Should(Equal("execute"))
		Expect(parsedMessage.MetaData.Finished).Should(Equal(false))
		Expect(len(parsedMessage.Data)).Should(Equal(2))
		Expect(parsedMessage.Data[0]["_raw"]).Should(Equal([]byte("you suck")))
		Expect(parsedMessage.Data[0]["_time"]).Should(Equal([]byte("1234567890")))
		Expect(parsedMessage.Data[0]["index"]).Should(Equal([]byte("main")))
		Expect(parsedMessage.Data[0]["sourcetype"]).Should(Equal([]byte("silly")))
		Expect(parsedMessage.Data[0]["source"]).Should(Equal([]byte("udp")))
		Expect(parsedMessage.Data[0]["host"]).Should(Equal([]byte("mrt.splunk.com")))

	})

	It("correctly parses simple CEC protocol 'getinfo message from splunk", func() {
		cecMessage := SplunkCecMessage{}
		parsedMessage, err := cecMessage.Parse(validGetinfoMessage)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(parsedMessage.ProtocolVersion).Should(Equal("chunked 1.0"))
		Expect(parsedMessage.MetadataLength).Should(Equal(23))
		Expect(parsedMessage.DataLength).Should(Equal(0))
		Expect(parsedMessage.MetaData.Action).Should(Equal("getinfo"))

	})

	It("correctly parses complex CEC protocol 'getinfo message from splunk", func() {
		cecMessage := SplunkCecMessage{}
		parsedMessage, err := cecMessage.Parse(validComplexGetinfoMessage)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(parsedMessage.ProtocolVersion).Should(Equal("chunked 1.0"))
		Expect(parsedMessage.MetadataLength).Should(Equal(444))
		Expect(parsedMessage.DataLength).Should(Equal(0))
		Expect(parsedMessage.MetaData.Action).Should(Equal("getinfo"))
		Expect(parsedMessage.MetaData.Preview).Should(Equal(false))
		Expect(parsedMessage.MetaData.Streaming_command_will_restart).Should(Equal(false))
		Expect(len(parsedMessage.MetaData.Searchinfo.Args)).Should(Equal(2))
		Expect(parsedMessage.MetaData.Searchinfo.Args[0]).Should(Equal("arg1"))
		Expect(len(parsedMessage.MetaData.Searchinfo.Raw_args)).Should(Equal(2))
		Expect(parsedMessage.MetaData.Searchinfo.Raw_args[0]).Should(Equal("rawarg1"))
		Expect(parsedMessage.MetaData.Searchinfo.Dispatch_dir).Should(Equal("/tmp"))
		Expect(parsedMessage.MetaData.Searchinfo.Sid).Should(Equal("unittest"))
		Expect(parsedMessage.MetaData.Searchinfo.App).Should(Equal("app"))
		Expect(parsedMessage.MetaData.Searchinfo.Username).Should(Equal("weberjas"))
		Expect(parsedMessage.MetaData.Searchinfo.Owner).Should(Equal("jasonw"))
		Expect(parsedMessage.MetaData.Searchinfo.Session_key).Should(Equal("123ABC"))
		Expect(parsedMessage.MetaData.Searchinfo.Splunkd_uri).Should(Equal("/custom"))
		Expect(parsedMessage.MetaData.Searchinfo.Splunk_version).Should(Equal("8.8.8"))
		Expect(parsedMessage.MetaData.Searchinfo.Search).Should(Equal("SPL String"))
		Expect(parsedMessage.MetaData.Searchinfo.Command).Should(Equal("commandName"))
		Expect(parsedMessage.MetaData.Searchinfo.Maxresultrows).Should(Equal(100))
		Expect(parsedMessage.MetaData.Searchinfo.Earliest_time).Should(Equal(123456))
		Expect(parsedMessage.MetaData.Searchinfo.Latest_time).Should(Equal(987654))

	})
})
