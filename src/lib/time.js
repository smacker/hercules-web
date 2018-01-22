import differenceInMonths from 'date-fns/difference_in_months';
import differenceInYears from 'date-fns/difference_in_years';

export function chooseDefaultResampling(begin, end) {
  begin = new Date(begin * 1000);
  end = new Date(end * 1000);

  const years = differenceInYears(end, begin);
  if (years >= 3) {
    return 'year';
  }
  const months = differenceInMonths(end, begin);
  if (months > 3 && months < 36) {
    return 'month';
  }
  return 'raw';
}
