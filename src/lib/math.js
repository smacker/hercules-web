// custom bundling
// import only what we need from mathjs

import mathCore from 'mathjs/core';
import bigNumberType from 'mathjs/lib/type/bignumber';
import matrixType from 'mathjs/lib/type/matrix';
import zerosFn from 'mathjs/lib/function/matrix/zeros';
import rangeFn from 'mathjs/lib/function/matrix/range';
import transposeFn from 'mathjs/lib/function/matrix/transpose';

const math = mathCore.create();
math.import(bigNumberType);
math.import(matrixType);
math.import(zerosFn);
math.import(rangeFn);
math.import(transposeFn);

export default math;
