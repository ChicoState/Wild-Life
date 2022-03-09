import {Buffer} from "buffer"

export function getBuffer(file: File) {
  return function(resolve){
    var reader = new FileReader();
    reader.readAsArrayBuffer(file);
    reader.onload = function() {
      var arrayBuffer = reader.result;
      var bytes: any
      if(typeof arrayBuffer === 'string') {
        bytes = new Uint8Array(arrayBuffer.length);
      } else {
        var raw = new Uint8Array(arrayBuffer);
        // var decoeder = new TextDecoder('utf-8');
        bytes = encodeBase64(raw);
      }
      resolve(bytes);
    }
  }
}

function encodeBase64(data: Uint8Array) {
   return Buffer.from(data).toString('base64');
}
