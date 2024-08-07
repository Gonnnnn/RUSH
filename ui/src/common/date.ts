const toYYYY년MM월DD일HH시MM분 = (date: Date) =>
  `${date.getFullYear()}년 ${date.getMonth() + 1}월 ${date.getDate()}일 ${date.getHours()}시 ${date.getMinutes()}분`;

const toYYslashMMslashDDspaceHHcolonMM = (date: Date) =>
  `${date.getFullYear()}/${twoDigit(date.getMonth() + 1)}/${twoDigit(date.getDate())} ${twoDigit(date.getHours())}:${twoDigit(date.getMinutes())}`;

const twoDigit = (n: number) => String(n).padStart(2, '0');

export { toYYYY년MM월DD일HH시MM분, toYYslashMMslashDDspaceHHcolonMM };
