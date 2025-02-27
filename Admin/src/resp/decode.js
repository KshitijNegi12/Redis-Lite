'use strict';

const parseRESP = (cmds) => {
	while(cmds.length) {
		const element = cmds.shift();
		switch (element[0]) {
			case '+':
                return [element.slice(1), cmds]
			case '*':
				const arrlen = element.slice(1);
				const arr = [];
				for (let j = 0; j < arrlen; j++) {
					const parsedContent = parseRESP(cmds);
					arr.push(parsedContent[0]);
					cmds = parsedContent[1];
				}
				return arr;
			case '$':
				const strlen = element.slice(1);
				const str = cmds.shift();
				return [str, cmds];
			case ':':
				const integer = element.slice(1);
				return [Number(integer), cmds];
			default:
				return [element, cmds];
		}
	}
};

module.exports = {parseRESP};