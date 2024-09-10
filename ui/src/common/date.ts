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

  const getOrdinalSuffix = (dayToConvert: number) => {
    const suffix = ['th', 'st', 'nd', 'rd'];
    if (dayToConvert > 10 && dayToConvert < 20) return dayToConvert + suffix[0];
    const suffixIndex = dayToConvert % 10;
    if (suffixIndex > 3) return dayToConvert + suffix[0];
    return dayToConvert + suffix[suffixIndex];
  };

  return `${month} ${getOrdinalSuffix(day)}`;
};

const twoDigit = (n: number) => String(n).padStart(2, '0');

export { toYYYY년MM월DD일HH시MM분, toYYslashMMslashDDspaceHHcolonMM, formatDateToMonthDate };
