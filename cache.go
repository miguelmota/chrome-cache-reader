package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	pack "github.com/miguelmota/chrome-cache-reader/pack"
)

func main() {
	var urls []string
	cacheDir := "~/Library/Caches/Google/Chrome/Profile 1/Cache/"
	outputDir := "./out"
	_ = outputDir
	_ = cacheDir

	cp := new(cacheParse)

	cp.parse(cacheDir, urls)
	//cache := cp.parse(cacheDir, urls)
	//cp.exportToHTML(cache, outputDir)
}

type cacheAddress struct {
}

// newCacheAddress ...
func newCacheAddress() *cacheAddress {
	return &cacheAddress{}
	/*
	   """
	   Object representing a Chrome Cache Address
	   """
	   SEPARATE_FILE = 0
	   RANKING_BLOCK = 1
	   BLOCK_256 = 2
	   BLOCK_1024 = 3
	   BLOCK_4096 = 4

	   typeArray = [("Separate file", 0),
	                ("Ranking block file", 36),
	                ("256 bytes block file", 256),
	                ("1k bytes block file", 1024),
	                ("4k bytes block file", 4096)]

	   def __init__(self, uint_32, path):
	       """
	       Parse the 32 bits of the uint_32
	       """
	       if uint_32 == 0:
	           raise CacheAddressError("Null Address")

	       #XXX Is self.binary useful ??
	       self.addr = uint_32
	       self.path = path

	       # Checking that the MSB is set
	       self.binary = bin(uint_32)
	       if len(self.binary) != 34:
	           raise CacheAddressError("Uninitialized Address")

	       self.blockType = int(self.binary[3:6], 2)

	       # If it is an address of a separate file
	       if self.blockType == CacheAddress.SEPARATE_FILE:
	           self.fileSelector = "f_%06x" % int(self.binary[6:], 2)
	       elif self.blockType == CacheAddress.RANKING_BLOCK:
	           self.fileSelector = "data_" + str(int(self.binary[10:18], 2))
	       else:
	           self.entrySize = CacheAddress.typeArray[self.blockType][1]
	           self.contiguousBlock = int(self.binary[8:10], 2)
	           self.fileSelector = "data_" + str(int(self.binary[10:18], 2))
	           self.blockNumber = int(self.binary[18:], 2)
	*/
}

func (s *cacheAddress) toString() {
	/*
	   string = hex(self.addr) + " ("
	   if self.blockType >= CacheAddress.BLOCK_256:
	       string += str(self.contiguousBlock) +\
	                 " contiguous blocks in "
	   string += CacheAddress.typeArray[self.blockType][0] +\
	             " : " + self.fileSelector + ")"
	   return string
	*/
}

// cacheBlock ...
// Object representing a block of the cache. It can be the index file or any
// other block type : 256B, 1024B, 4096B, Ranking Block.
// See /net/disk_cache/disk_format.h for details.
type cacheBlock struct {
	blockType int
	version   []byte
}

var INDEX_MAGIC = 0xC103CAC3
var BLOCK_MAGIC = 0xC104CAC3
var UNKNOWN = 0
var INDEX = 0
var BLOCK = 1

// newCacheBlock ...
// Parse the header of a cache file
func newCacheBlock(filename string) *cacheBlock {
	s := &cacheBlock{
		blockType: UNKNOWN,
	}

	fmt.Println(filename)
	header, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	b := []byte{195, 202, 3, 193}
	fmt.Println("YO", hex.EncodeToString(b))

	// Create BinaryPack object
	bp := new(pack.BinaryPack)
	var data [4]byte
	n, err := io.ReadFull(header, data[:])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("N", n)
	fmt.Println("data L", len(data[:]))
	fmt.Println("data", data[:])

	fmt.Println(hex.EncodeToString(data[:]))

	c, err := bp.CalcSize([]string{"I"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("calc size", c)
	magicU, err := bp.UnPack([]string{"I"}, data[:])
	if err != nil {
		log.Fatal(err)
	}

	magic, ok := magicU[0].(int)
	if !ok {
		log.Fatal("not ok")
	}
	fmt.Println("MAGIC L", len(magicU))
	fmt.Println("MAGIC", magic)
	fmt.Println("MAGIC LITTLE", binary.LittleEndian.Uint32(b))

	bs := []byte(strconv.Itoa(magic))
	fmt.Println(hex.EncodeToString(bs))

	_, err = header.Seek(2, 1)
	if err != nil {
		log.Fatal(err)
	}
	var data2 [2]byte
	_, err = io.ReadFull(header, data2[:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DATA2", data2)
	versionU, err := bp.UnPack([]string{"h"}, data2[:])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("VERSION L", len(versionU))
	version, ok := versionU[0].(int)
	if !ok {
		log.Fatal("NO OK")
	}

	fmt.Println("VERSION", version)
	/*
	   header = open(filename, 'rb')

	   # Read Magic Number
	   magic = struct.unpack('I', header.read(4))[0]
	   if magic == CacheBlock.BLOCK_MAGIC:
	       self.type = CacheBlock.BLOCK
	       header.seek(2, 1)
	       self.version = struct.unpack('h', header.read(2))[0]
	       self.header = struct.unpack('h', header.read(2))[0]
	       self.nextFile = struct.unpack('h', header.read(2))[0]
	       self.blockSize = struct.unpack('I', header.read(4))[0]
	       self.entryCount = struct.unpack('I', header.read(4))[0]
	       self.entryMax = struct.unpack('I', header.read(4))[0]
	       self.empty = []
	       for _ in range(4):
	           self.empty.append(struct.unpack('I', header.read(4))[0])
	       self.position = []
	       for _ in range(4):
	           self.position.append(struct.unpack('I', header.read(4))[0])
	   elif magic == CacheBlock.INDEX_MAGIC:
	       self.type = CacheBlock.INDEX
	       header.seek(2, 1)
	       self.version = struct.unpack('h', header.read(2))[0]
	       self.entryCount = struct.unpack('I', header.read(4))[0]
	       self.byteCount = struct.unpack('I', header.read(4))[0]
	       self.lastFileCreated = "f_%06x" % \
	                                  struct.unpack('I', header.read(4))[0]
	       header.seek(4*2, 1)
	       self.tableSize = struct.unpack('I', header.read(4))[0]
	   else:
	       header.close()
	       raise Exception("Invalid Chrome Cache File")
	   header.close()
	*/

	return s
}

func (s *cacheBlock) getType() int {
	return s.blockType
}

type cacheData struct {
}

// newCacheData ...
func newCacheData() *cacheData {
	return &cacheData{}
	/*
	   HTTP_HEADER = 0
	   UNKNOWN = 1

	   def __init__(self, address, size, isHTTPHeader=False):
	       """
	       It is a lazy evaluation object : the file is open only if it is
	       needed. It can parse the HTTP header if asked to do so.
	       See net/http/http_util.cc LocateStartOfStatusLine and
	       LocateEndOfHeaders for details.
	       """
	       self.size = size
	       self.address = address
	       self.type = CacheData.UNKNOWN

	       if isHTTPHeader and\
	          self.address.blockType != cacheAddress.CacheAddress.SEPARATE_FILE:
	           # Getting raw data
	           string = ""
	           block = open(self.address.path + self.address.fileSelector, 'rb')
	           block.seek(8192 + self.address.blockNumber*self.address.entrySize)
	           for _ in range(self.size):
	               string += struct.unpack('c', block.read(1))[0]
	           block.close()

	           # Finding the beginning of the request
	           start = re.search("HTTP", string)
	           if start == None:
	               return
	           else:
	               string = string[start.start():]

	           # Finding the end (some null characters : verified by experience)
	           end = re.search("\x00\x00", string)
	           if end == None:
	               return
	           else:
	               string = string[:end.end()-2]

	           # Creating the dictionary of headers
	           self.headers = {}
	           for line in string.split('\0'):
	               stripped = line.split(':')
	               self.headers[stripped[0].lower()] = \
	                   ':'.join(stripped[1:]).strip()
	           self.type = CacheData.HTTP_HEADER
	*/
}

func (s *cacheData) save() {
	/*
	   def save(self, filename=None):
	       """Save the data to the specified filename"""
	       if self.address.blockType == cacheAddress.CacheAddress.SEPARATE_FILE:
	           print(filename)
	           if os.path.isfile(filename):
	               shutil.copy(self.address.path + self.address.fileSelector,
	                       filename)
	       else:
	           output = open(filename, 'wb')
	           block = open(self.address.path + self.address.fileSelector, 'rb')
	           block.seek(8192 + self.address.blockNumber*self.address.entrySize)
	           output.write(block.read(self.size))
	           block.close()
	           output.close()
	*/
}

func (s *cacheData) data() {
	/*
	   """Returns a string representing the data"""
	   block = open(self.address.path + self.address.fileSelector, 'rb')
	   block.seek(8192 + self.address.blockNumber*self.address.entrySize)
	   data = block.read(self.size).decode('utf-8')
	   block.close()
	   return data
	*/
}

func (s *cacheData) toString() {
	/*
	   """
	   Display the type of cacheData
	   """
	   if self.type == CacheData.HTTP_HEADER:
	       if self.headers.has_key('content-type'):
	           return "HTTP Header %s" % self.headers['content-type']
	       else:
	           return "HTTP Header"
	   else:
	       return "Data"
	*/
}

type cacheEntry struct {
}

// newCacheEntry ...
func newCacheEntry() *cacheEntry {
	return &cacheEntry{}
	/*
	   """
	   See /net/disk_cache/disk_format.h for details.
	   """
	   STATE = ["Normal",
	            "Evicted (data were deleted)",
	            "Doomed (shit happened)"]

	   def __init__(self, address):
	       """
	       Parse a Chrome Cache Entry at the given address
	       """
	       self.httpHeader = None
	       block = open(address.path + address.fileSelector, 'rb')

	       # Going to the right entry
	       block.seek(8192 + address.blockNumber*address.entrySize)

	       # Parsing basic fields
	       self.hash = struct.unpack('I', block.read(4))[0]
	       self.next = struct.unpack('I', block.read(4))[0]
	       self.rankingNode = struct.unpack('I', block.read(4))[0]
	       self.usageCounter = struct.unpack('I', block.read(4))[0]
	       self.reuseCounter = struct.unpack('I', block.read(4))[0]
	       self.state = struct.unpack('I', block.read(4))[0]
	       self.creationTime = datetime.datetime(1601, 1, 1) + \
	                           datetime.timedelta(microseconds=\
	                               struct.unpack('Q', block.read(8))[0])
	       self.keyLength = struct.unpack('I', block.read(4))[0]
	       self.keyAddress = struct.unpack('I', block.read(4))[0]


	       dataSize = []
	       for _ in range(4):
	           dataSize.append(struct.unpack('I', block.read(4))[0])

	       self.data = []
	       for index in range(4):
	           addr = struct.unpack('I', block.read(4))[0]
	           try:
	               addr = cacheAddress.CacheAddress(addr, address.path)
	               self.data.append(cacheData.CacheData(addr, dataSize[index],
	                                                    True))
	           except cacheAddress.CacheAddressError:
	               pass

	       # Find the HTTP header if there is one
	       for data in self.data:
	           if data.type == cacheData.CacheData.HTTP_HEADER:
	               self.httpHeader = data
	               break

	       self.flags = struct.unpack('I', block.read(4))[0]

	       # Skipping pad
	       block.seek(5*4, 1)

	       # Reading local key
	       if self.keyAddress == 0:
	           self.key = block.read(self.keyLength).decode('ascii')
	       # Key stored elsewhere
	       else:
	           addr = cacheAddress.CacheAddress(self.keyAddress, address.path)

	           # It is probably an HTTP header
	           self.key = cacheData.CacheData(addr, self.keyLength, True)

	       block.close()
	*/
}

func (s *cacheEntry) keyToStr() {
	/*
	   def keyToStr(self):
	       """
	       Since the key can be a string or a CacheData object, this function is an
	       utility to display the content of the key whatever type is it.
	       """
	       if self.keyAddress == 0:
	           return self.key
	       else:
	           return self.key.data()

	   def __str__(self):
	       string = "Hash: 0x%08x" % self.hash + '\n'
	       if self.next != 0:
	           string += "Next: 0x%08x" % self.next + '\n'
	       string += "Usage Counter: %d" % self.usageCounter + '\n'\
	                 "Reuse Counter: %d" % self.reuseCounter + '\n'\
	                 "Creation Time: %s" % self.creationTime + '\n'
	       if self.keyAddress != 0:
	           string += "Key Address: 0x%08x" % self.keyAddress + '\n'
	       string += "Key: %s" % self.key + '\n'
	       if self.flags != 0:
	           string += "Flags: 0x%08x" % self.flags + '\n'
	       string += "State: %s" % CacheEntry.STATE[self.state]
	       for data in self.data:
	           string += "\nData (%d bytes) at 0x%08x : %s" % (data.size,
	                                                           data.address.addr,
	                                                           data)
	       return string
	*/
}

type cacheParse struct {
}

// Reads the whole cache and store the collected data in a table
// or find out if the given list of urls is in the cache. If yes it
// return a list of the corresponding entries.
func (s *cacheParse) parse(path string, urls []string) {
	// Verifying that the path end with / (What happen on windows?)
	path = NormalizePath(path)
	path = path + "/"
	cacheBlock := newCacheBlock(path + "index")

	// Checking type
	//if cacheBlock.type() != CacheBlock.INDEX {
	if cacheBlock.getType() != INDEX {
		//log.Fatal("Invalid Index File")
	}

	index, err := os.Open(path + "index")
	if err != nil {
		log.Fatal(err)
	}

	// Skipping Header
	index.Seek(92*4, 0)

	var cache []*cacheEntry
	_ = cache
	// If no url is specified, parse the whole cache
	if len(urls) == 0 {
		/*
					 for key, value := range(cacheBlock.tableSize) {
			           raw = struct.unpack('I', index.read(4))[0]
			           if raw != 0:
			               entry = CacheEntry(CacheAddress(raw, path=path))
			               // Checking if there is a next item in the bucket because
			               // such entries are not stored in the Index File so they will
			               // be ignored during iterative lookup in the hash table
			               while entry.next != 0 {
			                   cache = append(cache,, entry)
			                   entry = newCacheEntry(newCacheAddress(entry.next, path=path))
										}
			               cache = append(cache, entry)
							}
		*/
	}
	/*
	   else:
	       # Find the entry for each url
	       for url in urls:
	           # Compute the key and seeking to it
	           hash = SuperFastHash.superFastHash(url)
	           key = hash & (cacheBlock.tableSize - 1)
	           index.seek(92*4 + key*4)

	           addr = struct.unpack('I', index.read(4))[0]
	           # Checking if the address is initialized (i.e. used)
	           if addr & 0x80000000 == 0:
	               print >> sys.stderr, \
	                     "\033[32m%s\033[31m is not in the cache\033[0m" % url

	           # Follow the chained list in the bucket
	           else:
	               entry = CacheEntry(CacheAddress(addr, path=path))
	               while entry.hash != hash and entry.next != 0:
	                   entry = CacheEntry(CacheAddress(entry.next, path=path))
	               if entry.hash == hash:
	                   cache.append(entry)
	*/
	fmt.Println("ret")
	//return cache
}

func (s *cacheParse) exportToHTML(cache string, outputDir string) {
	/*
	   """
	   Export the cache in html
	   """

	   # Checking that the directory exists and is writable
	   if not os.path.exists(outpath):
	       os.makedirs(outpath)
	   outpath = os.path.abspath(outpath) + '/'

	   index = open(outpath + "index.html", 'w')
	   index.write("<UL>")

	   for entry in cache:
	       # Adding a link in the index
	       if entry.keyLength > 100:
	           entry_name = entry.keyToStr()[:100] + "..."
	       else:
	           entry_name = entry.keyToStr()
	       index.write('<LI><a href="%08x">%s</a></LI>'%(entry.hash, entry_name))
	       # We handle the special case where entry_name ends with a slash
	       page_basename = entry_name.split('/')[-2] if entry_name.endswith('/') else entry_name.split('/')[-1]

	       # Creating the entry page
	       page = open(outpath + "%08x"%entry.hash, 'w')
	       page.write("""<!DOCTYPE html>
	                     <html lang="en">
	                     <head>
	                     <meta charset="utf-8">
	                     </head>
	                     <body>""")

	       # Details of the entry
	       page.write("<b>Hash</b>: 0x%08x<br />"%entry.hash)
	       page.write("<b>Usage Counter</b>: %d<br />"%entry.usageCounter)
	       page.write("<b>Reuse Counter</b>: %d<br />"%entry.reuseCounter)
	       page.write("<b>Creation Time</b>: %s<br />"%entry.creationTime)
	       page.write("<b>Key</b>: %s<br>"%entry.keyToStr())
	       page.write("<b>State</b>: %s<br>"%CacheEntry.STATE[entry.state])

	       page.write("<hr>")
	       if len(entry.data) == 0:
	           page.write("No data associated with this entry :-(")
	       for i in range(len(entry.data)):
	           if entry.data[i].type == CacheData.UNKNOWN:
	               # Extracting data into a file
	               name = hex(entry.hash) + "_" + str(i)
	               entry.data[i].save(outpath + name)

	               if entry.httpHeader != None and \
	                  entry.httpHeader.headers.has_key('content-encoding') and\
	                  entry.httpHeader.headers['content-encoding'] == "gzip":
	                   # XXX Highly inefficient !!!!!
	                   try:
	                       input = gzip.open(outpath + name, 'rb')
	                       output = open(outpath + name + "u", 'w')
	                       output.write(input.read())
	                       input.close()
	                       output.close()
	                       page.write('<a href="%su">%s</a>'%(name, page_basename))
	                   except IOError:
	                       page.write("Something wrong happened while unzipping")
	               else:
	                   page.write('<a href="%s">%s</a>'%(name ,
	                              entry.keyToStr().split('/')[-1]))


	               # If it is a picture, display it
	               if entry.httpHeader != None:
	                   if entry.httpHeader.headers.has_key('content-type') and\
	                      "image" in entry.httpHeader.headers['content-type']:
	                       page.write('<br /><img src="%s">'%(name))
	           # HTTP Header
	           else:
	               page.write("<u>HTTP Header</u><br />")
	               for key, value in entry.data[i].headers.items():
	                   page.write("<b>%s</b>: %s<br />"%(key, value))
	           page.write("<hr>")
	       page.write("</body></html>")
	       page.close()

	   index.write("</UL>")
	   index.close()
	*/
}

// UserHomeDir ...
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}

// NormalizePath ...
func NormalizePath(path string) string {
	// expand tilde
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(UserHomeDir(), path[2:])
	}

	return path
}
