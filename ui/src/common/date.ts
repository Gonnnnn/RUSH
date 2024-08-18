const toYYYY년MM월DD일HH시MM분 = (date: Date) =>
  `${date.getFullYear()}년 ${date.getMonth() + 1}월 ${date.getDate()}일 ${date.getHours()}시 ${date.getMinutes()}분`;

const toYYslashMMslashDDspaceHHcolonMM = (date: Date) =>
  `${date.getFullYear()}/${twoDigit(date.getMonth() + 1)}/${twoDigit(date.getDate())} ${twoDigit(date.getHours())}:${twoDigit(date.getMinutes())}`;

const formatDateToMonthDate = (date: Date) => {
  const months = [
    'January',
    'February',
    'March',
    'April',
    'May',
    'June',
    'July',
    'August',
    'September',
    'October',
    'November',
    'December',
  ];

  const day = date.getDate();
  const month = months[date.getMonth()];

  const getOrdinalSuffix = (datyToConvert: number) => {
    const suffix = ['th', 'st', 'nd', 'rd'];
    const value = datyToConvert % 100;
    if (value > 10 && value < 20) return datyToConvert + suffix[0];
    return datyToConvert + suffix[value % 10];
  };

  return `${month} ${getOrdinalSuffix(day)}`;
};

const twoDigit = (n: number) => String(n).padStart(2, '0');

export { toYYYY년MM월DD일HH시MM분, toYYslashMMslashDDspaceHHcolonMM, formatDateToMonthDate };
