'use strict';

const toRESP = (obj) =>{
	let resp = '';
	switch (typeof obj) {
		case 'object':
			if (obj.constructor === Array) {
				const arrLen = obj.length;
				resp += `*${arrLen}\r\n`;
				for (let i = 0; i < arrLen; i++) {
					resp += toRESP(obj[i]);
				}
			}
			return resp;
		case 'string':
            return `$${obj.length}\r\n${obj}\r\n`;
		case 'number':
			return `:${obj}\r\n`;
		default:
			break;
	}
}

module.exports = {toRESP};