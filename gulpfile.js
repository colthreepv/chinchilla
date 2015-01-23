// gulp task requires
var gulp = require('gulp');
var gutil = require('gulp-util');
var sequence = require('run-sequence');
var nunjucks = require('gulp-nunjucks-html');
var bowerFiles = require('main-bower-files');
var del = require('del');
var bower = require('bower');

// stdlib
var path = require('path');

// make this go away someday
var gdebug = require('gulp-debug');

gulp.task('clean', function (done) {
  del(['build/'], done);
});

gulp.task('install-libs', function (done) {
  bower.commands.install([], { save: true }, { interactive: false })
    .on('end', function (installed) {
      done();
    });
});

gulp.task('copy-libs', function () {
  return gulp.src(bowerFiles({ env: 'development' }))
    .pipe(gulp.dest('build/libs/'));
});

/**
 * filters libraries to be included in template
 * Ideal for Array.prototype.filter function.
 * @param  {string} ext: '.js' || '.css' and so on
 * @param  {string} skipThose: library to filter out
 * @param  {array}  skipThose: array of libraries with such extension to filter out
 * @return {boolean} true to include the file, false to exclude it.
 */
function filterByExt (ext, skipThose) {
  return function (file) {
    if (Array.isArray(skipThose)) {
      if (skipThose.indexOf(path.basename(file, ext)) >= 0) { // if an occurence is found
        return false;
      }
    }
    // if that's not array, it's a string
    if (!!skipThose) {
      if (path.basename(file, ext) === skipThose) {
        return false;
      }
    }
    if (path.extname(file) === ext) {
      return true;
    } else {
      return false;
    }
  };
}

gulp.task('compile-html', function (done) {
  var files = bowerFiles({ env: 'development' });
  var jsFiles = files.filter(filterByExt('.js', ['react', 'bootstrap'])).map(function (l) { return 'libs/' + path.basename(l); });
  var cssFiles = files.filter(filterByExt('.css')).map(function (l) { return 'libs/' + path.basename(l); });
  return gulp.src('index.html')
    .pipe(nunjucks({
      locals: {
        build: 'development',
        jsFiles: jsFiles,
        cssFiles: cssFiles
      }
    }))
    .on('error', gutil.log)
    .pipe(gulp.dest('build/'));
});

gulp.task('watch', function () {
  gulp.watch(['src/*.js', 'src/**/*.js'], ['compile-src']);
  gulp.watch(['index.html'], ['compile-html']);
});

gulp.task('default', function (done) {
  sequence('clean', 'install-libs', ['compile-html', 'copy-libs'], 'watch');
});
