define([
  'lodash',
],
function (_) {
  "use strict";

  return function Settings (options) {
    var defaults = {
      datasources                   : {},
      window_title_prefix           : 'Grafana - ',
      panels                        : {},
      new_panel_title: '未命名',
      playlist_timespan: "1m",
      unsaved_changes_warning: true,
      appSubUrl: ""
    };

    return _.extend({}, defaults, options);
  };
});
