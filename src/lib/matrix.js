import math from 'mathjs';
import isEqual from 'date-fns/is_equal';
import differenceInDays from 'date-fns/difference_in_days';

export function toMatrix(data) {
  const matrix = math.zeros(data.length, data.length);
  data.forEach((row, i) => {
    row.forEach((cell, j) => {
      matrix.set([i, j], cell);
    });
  });
  return matrix;
}

export function interpolate(matrix, granularity, sampling) {
  const shape = matrix.size();
  const daily = math.zeros(shape[0] * granularity, shape[1] * sampling);

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
        const k = matrix.get([y, x]) / startVal;
        const scale = (x + 1) * sampling - startIndex;
        math.range(y * granularity, (y + 1) * granularity).forEach(i => {
          const initial = daily.get([i, startIndex - 1]);
          math.range(startIndex, (x + 1) * sampling).forEach(j => {
            const val = initial * (1 + (k - 1) * (j - startIndex + 1) / scale);
            daily.set([i, j], val);
          });
        });
      };

      const grow = (finishIndex, finishVal) => {
        const initial = x > 0 ? matrix.get([y, x - 1]) : 0;
        let startIndex = x * sampling;
        if (startIndex < y * granularity) {
          startIndex = y * granularity;
        }
        if (finishIndex == startIndex) {
          return;
        }
        const avg = (finishVal - initial) / (finishIndex - startIndex);
        math.range(x * sampling, finishIndex).forEach(j => {
          math.range(startIndex, j + 1).forEach(i => {
            daily.set([i, j], avg);
          });
        });
        math.range(x * sampling, finishIndex).forEach(j => {
          math.range(y * granularity, x * sampling).forEach(i => {
            daily.set([i, j], daily.get([i, j - 1]));
          });
        });
      };

      if ((y + 1) * granularity >= (x + 1) * sampling) {
        if (y * granularity <= x * sampling) {
          grow((x + 1) * sampling, matrix.get([y, x]));
        } else if ((x + 1) * sampling > y * granularity) {
          grow((x + 1) * sampling, matrix.get([y, x]));
          // FIXME figure out why we need it
          const avg =
            matrix.get([y, x]) / ((x + 1) * sampling - y * granularity);
          math.range(y * granularity, (x + 1) * sampling).forEach(j => {
            math.range(y * granularity, j + 1).forEach(i => {
              daily.set([i, j], avg);
            });
          });
        }
      } else if ((y + 1) * granularity >= x * sampling) {
        const v1 = matrix.get([y, x - 1]);
        const v2 = matrix.get([y, x]);
        const delta = (y + 1) * granularity - x * sampling;
        let previous = 0;
        let scale;
        if (x > 0 && (x - 1) * sampling >= y * granularity) {
          if (x > 1) {
            previous = matrix.get([y, x - 2]);
          }
          scale = sampling;
        } else {
          scale = x == 0 ? sampling : x * sampling - y * granularity;
        }
        let peak = v1 + (v1 - previous) / scale * delta;
        if (v2 > peak) {
          // we need to adjust the peak, it may not be less than the decayed value
          if (x < shape[1] - 1) {
            const k = (v2 - matrix.get([y, x + 1])) / sampling;
            peak =
              matrix.get([y, x]) +
              k * ((x + 1) * sampling - (y + 1) * granularity);
          } else {
            peak = v2;
          }
        }
        grow((y + 1) * granularity, peak);
        decay((y + 1) * granularity, peak);
      } else {
        decay(x * sampling, matrix.get([y, x - 1]));
      }
    }
  }

  return daily;
}

function aggregate(datesRange, { begin, end, data }, granularity, sampling) {
  const daily = interpolate(
    math.transpose(toMatrix(data)),
    granularity,
    sampling
  );
  const dailyShape = daily.size();

  const matrix = math.zeros(datesRange.length - 1, dailyShape[1]);

  datesRange.forEach((d, i) => {
    if (i === 0) {
      return;
    }

    const istart = differenceInDays(datesRange[i - 1], begin);
    const ifinish = differenceInDays(d, begin);

    const x = math.range(0, dailyShape[1]);
    let val = daily.subset(math.index(istart, x));
    for (let z = istart + 1; z < ifinish; z++) {
      val = math.add(val, daily.subset(math.index(z, x)));
    }
    matrix.subset(math.index(i - 1, math.range(0, dailyShape[1])), val);
  });
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
    data: math.transpose(matrix).toArray(),
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
    data: math.transpose(matrix).toArray(),
    keys: datesRange.slice(0, datesRange.length - 1).map(i => i.getFullYear())
  };
}
