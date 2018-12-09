import isEqual from 'date-fns/is_equal';
import differenceInDays from 'date-fns/difference_in_days';

export function transpose(matrix) {
  const rows = matrix.length;
  const cols = matrix[0].length;
  const newMatrix = [];

  for (let j = 0; j < cols; j++) {
    newMatrix[j] = Array(rows);
  }

  for (let i = 0; i < rows; i++) {
    for (let j = 0; j < cols; j++) {
      newMatrix[j][i] = matrix[i][j];
    }
  }

  return newMatrix;
}

export function sumByColumn(matrix) {
  const val = [];
  for (let i = 0; i < matrix.length; i++) {
    for (let j = 0; j < matrix[i].length; j++) {
      val[i] = (val[i] || 0) + matrix[i][j];
    }
  }
  return val;
}

export function interpolate(matrix, granularity, sampling) {
  const shape = [matrix.length, matrix[0].length];
  const daily = Array(shape[0] * granularity);
  for (let i = 0; i < daily.length; i++) {
    daily[i] = Array(shape[1] * sampling).fill(0);
  }

  for (let y = 0; y < shape[0]; y++) {
    for (let x = 0; x < shape[1]; x++) {
      // the future is zeros
      if (y * granularity > (x + 1) * sampling) {
        continue;
      }

      const decay = (startIndex, startVal) => {
        if (startVal === 0) {
          return;
        }
        const k = matrix[y][x] / startVal;
        const scale = (x + 1) * sampling - startIndex;
        for (let i = y * granularity; i < (y + 1) * granularity; i++) {
          const initial = daily[i][startIndex - 1];
          for (let j = startIndex; j < (x + 1) * sampling; j++) {
            const val =
              initial * (1 + ((k - 1) * (j - startIndex + 1)) / scale);
            daily[i][j] = val;
          }
        }
      };

      const grow = (finishIndex, finishVal) => {
        const initial = x > 0 ? matrix[y][x - 1] : 0;
        let startIndex = x * sampling;
        if (startIndex < y * granularity) {
          startIndex = y * granularity;
        }
        if (finishIndex == startIndex) {
          return;
        }
        const avg = (finishVal - initial) / (finishIndex - startIndex);
        for (let j = x * sampling; j < finishIndex; j++) {
          for (let i = startIndex; i < j + 1; i++) {
            daily[i][j] = avg;
          }
        }
        for (let j = x * sampling; j < finishIndex; j++) {
          for (let i = y * granularity; i < x * sampling; i++) {
            daily[i][j] = daily[i][j - 1];
          }
        }
      };

      if ((y + 1) * granularity >= (x + 1) * sampling) {
        if (y * granularity <= x * sampling) {
          grow((x + 1) * sampling, matrix[y][x]);
        } else if ((x + 1) * sampling > y * granularity) {
          grow((x + 1) * sampling, matrix[y][x]);
          // FIXME figure out why we need it
          const avg = matrix[y][x] / ((x + 1) * sampling - y * granularity);
          for (let j = y * granularity; j < (x + 1) * sampling; j++) {
            for (let i = y * granularity; i < j + 1; i++) {
              daily[i][j] = avg;
            }
          }
        }
      } else if ((y + 1) * granularity >= x * sampling) {
        const v1 = matrix[y][x - 1];
        const v2 = matrix[y][x];
        const delta = (y + 1) * granularity - x * sampling;
        let previous = 0;
        let scale;
        if (x > 0 && (x - 1) * sampling >= y * granularity) {
          if (x > 1) {
            previous = matrix[y][x - 2];
          }
          scale = sampling;
        } else {
          scale = x == 0 ? sampling : x * sampling - y * granularity;
        }
        let peak = v1 + ((v1 - previous) / scale) * delta;
        if (v2 > peak) {
          // we need to adjust the peak, it may not be less than the decayed value
          if (x < shape[1] - 1) {
            const k = (v2 - matrix[y][x + 1]) / sampling;
            peak =
              matrix[y][x] + k * ((x + 1) * sampling - (y + 1) * granularity);
          } else {
            peak = v2;
          }
        }
        grow((y + 1) * granularity, peak);
        decay((y + 1) * granularity, peak);
      } else {
        decay(x * sampling, matrix[y][x - 1]);
      }
    }
  }

  return daily;
}

function aggregate(datesRange, { begin, data }, granularity, sampling) {
  const daily = interpolate(transpose(data), granularity, sampling);
  const dailyShape = [daily.length, daily[0].length];

  const matrix = Array(datesRange.length - 1);

  for (let i = 1; i < datesRange.length; i++) {
    const istart = differenceInDays(datesRange[i - 1], begin);
    const ifinish = differenceInDays(datesRange[i], begin);

    const val = Array(dailyShape[1]);
    for (let z = istart; z < ifinish; z++) {
      for (let j = 0; j < dailyShape[1]; j++) {
        val[j] = (val[j] || 0) + daily[z][j];
      }
    }

    matrix[i - 1] = val;
  }
  return matrix;
}

const monthNames = {
  0: 'Jan',
  1: 'Feb',
  2: 'Mar',
  3: 'Apr',
  4: 'May',
  5: 'Jun',
  6: 'Jul',
  7: 'Aug',
  8: 'Sep',
  9: 'Oct',
  10: 'Nov',
  11: 'Dec'
};

export function toMonths(
  { begin, end, data },
  granularity = 30,
  sampling = 30
) {
  begin = new Date(begin * 1000);
  end = new Date(end * 1000);

  // array of dates like:
  // [24 Apr 2017, 1 Mar 2017, 1 Apr 2017, 28 Apr 2017]
  const datesRange = [begin];
  const stop = new Date(end.getFullYear(), end.getMonth(), 1);
  let current = begin;
  for (let month = begin.getMonth() + 1; !isEqual(current, stop); month++) {
    current = new Date(begin.getFullYear(), month, 1);
    datesRange.push(current);
  }
  if (!isEqual(datesRange[datesRange.length - 1], end)) {
    datesRange.push(end);
  }

  const matrix = aggregate(
    datesRange,
    { begin, end, data },
    granularity,
    sampling
  );

  return {
    data: transpose(matrix),
    keys: datesRange
      .slice(0, datesRange.length - 1)
      .map(i => monthNames[i.getMonth()] + ' ' + i.getFullYear())
  };
}

export function toYears({ begin, end, data }, granularity = 30, sampling = 30) {
  begin = new Date(begin * 1000);
  end = new Date(end * 1000);

  // array of dates like:
  // [24 Apr 2015, 1 Jan 2016, 1 Jan 2017, 1 Jan 2018, 30 Mar 2018]
  const datesRange = [begin];
  for (let year = begin.getFullYear() + 1; year <= end.getFullYear(); year++) {
    datesRange.push(new Date(year, 0, 1));
  }
  if (!isEqual(datesRange[datesRange.length - 1], end)) {
    datesRange.push(end);
  }

  const matrix = aggregate(
    datesRange,
    { begin, end, data },
    granularity,
    sampling
  );

  return {
    data: transpose(matrix),
    keys: datesRange.slice(0, datesRange.length - 1).map(i => i.getFullYear())
  };
}
