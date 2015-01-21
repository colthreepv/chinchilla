// gulp task requires
var gulp = require('gulp');
var sequence = require('run-sequence');
var nunjucks = require('gulp-nunjucks-html');

// stdlib
var path = require('path');

// globals for this gulpfile
var bowerDir = 'bower_components';
var jsFilesDev = [
  path.join(bowerDir, 'cash/build/debug/') + 'cash.js',
  path.join(bowerDir, 'delorean/dist/') + 'delorean.js'
];
var jsFilesDist = [];

gulp.task('copy-libs', function () {
  gulp.src(jsFilesDev)
    .pipe(gulp.dest('build/libs/'));
});

gulp.task('compile-html', function (done) {
  return gulp.src('index.html')
    .pipe(nunjucks({
      locals: {
        build: 'development',
        jsFiles: jsFilesDev.map(function (j) { return 'libs/' + path.basename(j); }),
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
  sequence(['compile-html', 'copy-libs'], 'watch');
});
