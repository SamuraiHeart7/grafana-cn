define([
  'angular',
],
function (angular) {
  'use strict';

  var module = angular.module('grafana.controllers');

  module.controller('AdminListOrgsCtrl', function($scope, backendSrv) {

    $scope.init = function() {
      $scope.getOrgs();
    };

    $scope.getOrgs = function() {
      backendSrv.get('/api/orgs').then(function(orgs) {
        $scope.orgs = orgs;
      });
    };

    $scope.deleteOrg = function(org) {
      $scope.appEvent('confirm-modal', {
        title: '删除',
        text: '确定要删除部门 ' + org.name + '?',
        text2: '此部门下所有的仪表盘都会被删除！',
        icon: 'fa-trash',
        yesText: '删除',
        onConfirm: function() {
          backendSrv.delete('/api/orgs/' + org.id).then(function() {
            $scope.getOrgs();
          });
        }
      });
    };

    $scope.init();

  });

});
