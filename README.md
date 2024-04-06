# zebug

This is a first pass at a tool for debugging gzip archives.
It prints gzip headers, gzip trailers, and DEFLATE block metadata.

The output is jsonlines where DEFLATE blocks are indented slightly to show the nesting.

```console
$ echo "My name is Jon Johnson I don't come from Wisconsin." | gzip | zebug
{"in":10,"header":{"modtime":"2024-04-06T15:19:02-07:00","os":3}}
  {"type":"01","in":58,"out":52,"final":true}
{"in":58,"out":52,"trailer":{"crc32":1401863900,"isize":52}}
```

## Examples

### Debugging `apk-tools`

This _would_ have been particularly useful when debugging a recently introduced issue in `apk-tools`.


Most APKs are exactly 3 gzip members for the signature, control, and data sections:

```console
$ curl -sL https://dl-cdn.alpinelinux.org/alpine/edge/main/aarch64/autoconf-archive-2023.02.20-r0.apk | zebug
{"in":10,"header":{"os":3}}
  {"type":"01","in":656,"out":1024,"final":true}
{"in":656,"out":1024,"trailer":{"crc32":148590525,"isize":1024}}
{"in":674,"out":1024,"header":{"os":3}}
  {"type":"10","in":1226,"out":3584,"final":true}
{"in":1226,"out":3584,"trailer":{"crc32":287079271,"isize":2560}}
{"in":1244,"out":3584,"header":{"os":3}}
  {"type":"10","in":24225,"out":134656}
  {"type":"01","in":24226,"out":134656}
  {"type":"01","in":24227,"out":134656}
  {"type":"10","in":37238,"out":265728}
  {"type":"01","in":37239,"out":265728}
  {"type":"01","in":37240,"out":265728}
  {"type":"01","in":37241,"out":265728}
  {"type":"10","in":56029,"out":396800}
  {"type":"00","in":56033,"out":396800}
  {"type":"10","in":66659,"out":527872}
  {"type":"00","in":66663,"out":527872}
  {"type":"10","in":80717,"out":658944}
  {"type":"00","in":80722,"out":658944}
  {"type":"10","in":101997,"out":790016}
  {"type":"00","in":102001,"out":790016}
  {"type":"10","in":126357,"out":921088}
  {"type":"00","in":126362,"out":921088}
  {"type":"10","in":148334,"out":1052160}
  {"type":"10","in":169780,"out":1183232}
  {"type":"10","in":181689,"out":1314304}
  {"type":"01","in":181690,"out":1314304}
  {"type":"01","in":181691,"out":1314304}
  {"type":"01","in":181692,"out":1314304}
  {"type":"10","in":194429,"out":1445376}
  {"type":"01","in":194430,"out":1445376}
  {"type":"01","in":194431,"out":1445376}
  {"type":"10","in":220154,"out":1543680}
  {"type":"10","in":220895,"out":1576448}
  {"type":"00","in":220899,"out":1576448}
  {"type":"10","in":241498,"out":1707520}
  {"type":"10","in":263493,"out":1838592}
  {"type":"00","in":263497,"out":1838592}
  {"type":"10","in":287346,"out":1969664}
  {"type":"00","in":287351,"out":1969664}
  {"type":"10","in":306630,"out":2100736}
  {"type":"00","in":306634,"out":2100736}
  {"type":"10","in":332636,"out":2231808}
  {"type":"01","in":332637,"out":2231808}
  {"type":"01","in":332638,"out":2231808}
  {"type":"01","in":332639,"out":2231808}
  {"type":"10","in":358926,"out":2362880}
  {"type":"00","in":358931,"out":2362880}
  {"type":"10","in":379321,"out":2493952}
  {"type":"01","in":379322,"out":2493952}
  {"type":"01","in":379323,"out":2493952}
  {"type":"01","in":379324,"out":2493952}
  {"type":"10","in":394591,"out":2625024}
  {"type":"00","in":394595,"out":2625024}
  {"type":"10","in":415968,"out":2756096}
  {"type":"10","in":438822,"out":2887168}
  {"type":"00","in":438826,"out":2887168}
  {"type":"10","in":460469,"out":3018240}
  {"type":"00","in":460473,"out":3018240}
  {"type":"10","in":482590,"out":3149312}
  {"type":"01","in":482591,"out":3149312}
  {"type":"10","in":488018,"out":3188224,"final":true}
{"in":488018,"out":3188224,"trailer":{"crc32":2277435338,"isize":3184640}}
```

This APK was built with [`pargzip`](https://pkg.go.dev/golang.org/x/build/pargzip), so its data section is split up into multiple gzip members:

```console
$ curl -sL https://packages.wolfi.dev/os/aarch64/autoconf-archive-2023.02.20-r0.apk | zebug
{"in":10}
  {"type":"01","in":649}
  {"type":"00","in":654,"final":true}
{"in":654,"out":1024,"trailer":{"crc32":2691215862,"isize":1024}}
{"in":672,"out":1024}
  {"type":"10","in":1004,"out":1024}
  {"type":"00","in":1009,"out":1024,"final":true}
{"in":1009,"out":2048,"trailer":{"crc32":2446349300,"isize":1024}}
{"in":1027,"out":2048}
  {"type":"10","in":22340,"out":100352}
  {"type":"10","in":47254,"out":329728}
  {"type":"10","in":69171,"out":591872}
  {"type":"10","in":93356,"out":722944}
  {"type":"10","in":116189,"out":886784}
  {"type":"10","in":139090,"out":1017856}
  {"type":"10","in":144151,"out":1050624}
  {"type":"00","in":144155,"out":1050624,"final":true}
{"in":144155,"out":1050624,"trailer":{"crc32":808686515,"isize":1048576}}
{"in":144173,"out":1050624}
  {"type":"10","in":166426,"out":1181696}
  {"type":"10","in":189408,"out":1443840}
  {"type":"10","in":211365,"out":1542144}
  {"type":"10","in":233654,"out":1673216}
  {"type":"10","in":257258,"out":1837056}
  {"type":"10","in":281818,"out":1968128}
  {"type":"10","in":299111,"out":2099200}
  {"type":"00","in":299115,"out":2099200,"final":true}
{"in":299115,"out":2099200,"trailer":{"crc32":2144299853,"isize":1048576}}
{"in":299133,"out":2099200}
  {"type":"10","in":321874,"out":2197504}
  {"type":"10","in":346152,"out":2328576}
  {"type":"10","in":369496,"out":2459648}
  {"type":"10","in":392601,"out":2656256}
  {"type":"10","in":415825,"out":2787328}
  {"type":"10","in":438898,"out":2918400}
  {"type":"10","in":461770,"out":3049472}
  {"type":"10","in":472259,"out":3147776}
  {"type":"00","in":472264,"out":3147776,"final":true}
{"in":472264,"out":3147776,"trailer":{"crc32":4046799306,"isize":1048576}}
{"in":472282,"out":3147776}
  {"type":"10","in":493468,"out":3213312}
  {"type":"10","in":515778,"out":3278848}
  {"type":"10","in":539255,"out":3508224}
  {"type":"10","in":562578,"out":3639296}
  {"type":"10","in":586047,"out":3803136}
  {"type":"10","in":610367,"out":3966976}
  {"type":"10","in":634362,"out":4098048}
  {"type":"10","in":643208,"out":4196352}
  {"type":"00","in":643213,"out":4196352,"final":true}
{"in":643213,"out":4196352,"trailer":{"crc32":4138702688,"isize":1048576}}
{"in":643231,"out":4196352}
  {"type":"10","in":665460,"out":4294656}
  {"type":"10","in":686700,"out":4392960}
  {"type":"10","in":697705,"out":4425728}
  {"type":"10","in":708567,"out":4491264}
  {"type":"10","in":719690,"out":4556800}
  {"type":"10","in":730790,"out":4589568}
  {"type":"10","in":741737,"out":4655104}
  {"type":"10","in":752719,"out":4687872}
  {"type":"10","in":763618,"out":4753408}
  {"type":"10","in":779033,"out":4884480}
  {"type":"10","in":781872,"out":4950016}
  {"type":"00","in":781876,"out":4950016,"final":true}
{"in":781876,"out":4975616,"trailer":{"crc32":1065729148,"isize":779264}}
```

The fourth gzip member starts at offset 1050624, which is in the middle of a file that `apk-tools` was reporting as a BAD archive.
I had to manually do some janky things with `tail` and `xxd` to figure this out, so I built this in case I ever need to do anything like that again.

### BTYPE Distribution of DEFLATE Blocks

It might be interesting to know how many blocks have a compressed (01 or 10) vs uncompressed (00) BTYPE.

```console
$ zebug < ubuntu.tar.gz | jq .type | grep -v null | sort | uniq -c
  33 "00"
1297 "10"
```

And also how much of a given gzip archive is made up of uncompressed (BTYPE 00) data.

```console
$ zebug < ubuntu.tar.gz | jq .Size | grep -v null | paste -s -d+ - | bc
541285
```

## TODO

* Web mode.
* Further breakdown of DEFLATE blocks.
    * Huffman trees
    * Pointers
    * Window utilization
