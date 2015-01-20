// gulp task requires
var gulp = require('gulp');
var browserify = require('browserify');
var namedStream = require('vinyl-source-stream');
var sequence = require('run-sequence');
var nunjucks = require('gulp-nunjucks-html');

// stdlib
var path = require('path');

gulp.task('compile-src', function () {
  var b = browserify(path.join(__dirname, 'src/') + 'main.js', {
    entry: true,
    debug: true
  });

  return b.bundle()
    .pipe(namedStream('app.js'))
    .pipe(gulp.dest('build/'));
});

gulp.task('compile-html', function (done) {
  return gulp.src('index.html')
    .pipe(nunjucks({
      locals: {
        jsFile: 'app.js',
        cssFile: 'not-available.css'
      }
    }))
    .pipe(gulp.dest('build/'));
});

gulp.task('watch', function () {
  gulp.watch(['src/*.js', 'src/**/*.js'], ['compile-src']);
  gulp.watch(['index.html'], ['compile-html']);
});

gulp.task('default', function (done) {
  sequence(['compile-src', 'compile-html'], 'watch');
});
