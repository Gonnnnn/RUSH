import { describe, expect, it } from 'vitest';
import { toYYYY년MM월DD일HH시MM분, toYYslashMMslashDDspaceHHcolonMM, formatDateToMonthDate } from './date';

describe('toYYYY년MM월DD일HH시MM분', () => {
  it('Success cases', () => {
    expect(toYYYY년MM월DD일HH시MM분(new Date('2024-01-01T00:12:00'))).toBe('2024년 1월 1일 0시 12분');
    expect(toYYYY년MM월DD일HH시MM분(new Date('2024-02-28T18:00:00'))).toBe('2024년 2월 28일 18시 0분');
    expect(toYYYY년MM월DD일HH시MM분(new Date('2024-03-31T01:12:00'))).toBe('2024년 3월 31일 1시 12분');
    expect(toYYYY년MM월DD일HH시MM분(new Date('2024-04-30T11:00:00'))).toBe('2024년 4월 30일 11시 0분');
  });
});

describe('toYYslashMMslashDDspaceHHcolonMM', () => {
  it('Success cases', () => {
    expect(toYYslashMMslashDDspaceHHcolonMM(new Date('2024-01-01T00:12:00'))).toBe('2024/01/01 00:12');
    expect(toYYslashMMslashDDspaceHHcolonMM(new Date('2024-02-28T00:00:00'))).toBe('2024/02/28 00:00');
    expect(toYYslashMMslashDDspaceHHcolonMM(new Date('2024-03-31T00:01:12'))).toBe('2024/03/31 00:01');
    expect(toYYslashMMslashDDspaceHHcolonMM(new Date('2024-04-30T00:00:00'))).toBe('2024/04/30 00:00');
  });
});

describe('formatDateToMonthDate', () => {
  it('Success cases', () => {
    expect(formatDateToMonthDate(new Date('2024-01-01T00:12:00'))).toBe('January 1st');
    expect(formatDateToMonthDate(new Date('2024-02-02T00:12:00'))).toBe('February 2nd');
    expect(formatDateToMonthDate(new Date('2024-03-03T00:12:00'))).toBe('March 3rd');
    expect(formatDateToMonthDate(new Date('2024-04-04T00:12:00'))).toBe('April 4th');
    expect(formatDateToMonthDate(new Date('2024-05-11T00:12:00'))).toBe('May 11th');
    expect(formatDateToMonthDate(new Date('2024-06-12T00:12:00'))).toBe('June 12th');
    expect(formatDateToMonthDate(new Date('2024-07-21T00:12:00'))).toBe('July 21st');
    expect(formatDateToMonthDate(new Date('2024-08-22T00:12:00'))).toBe('August 22nd');
    expect(formatDateToMonthDate(new Date('2024-09-23T00:12:00'))).toBe('September 23rd');
    expect(formatDateToMonthDate(new Date('2024-10-24T00:12:00'))).toBe('October 24th');
    expect(formatDateToMonthDate(new Date('2024-11-30T00:12:00'))).toBe('November 30th');
    expect(formatDateToMonthDate(new Date('2024-12-31T00:12:00'))).toBe('December 31st');
  });
});
