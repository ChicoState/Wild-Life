import {Buffer} from "buffer"

export function getBuffer(file: File) {
  return function(resolve:any){
    let reader = new FileReader();
    reader.readAsArrayBuffer(file);
    reader.onload = function() {
      let arrayBuffer = reader.result;
      let bytes: any
      // could be a string
      if(typeof arrayBuffer === 'string') {
        bytes = new Uint8Array(arrayBuffer.length);
      }else {
        // need to check if arrayBuffer is null
        let t_buff = (arrayBuffer !== null) ? arrayBuffer : new ArrayBuffer(0);
        let raw = new Uint8Array(t_buff);
        bytes = encodeBase64(raw);
      }
      resolve(bytes);
    }
  }
}

//turns a Uint8Array into a base64 string
function encodeBase64(data: Uint8Array) {
   return Buffer.from(data).toString('base64');
}
