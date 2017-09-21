angular.module('swarm', [])
.controller('SwarmController', ['$q', '$scope', 'SystemService', 'NodeService', 'Pagination', 'Notifications', 'StateManager', 'Authentication',
function ($q, $scope, SystemService, NodeService, Pagination, Notifications, StateManager, Authentication) {
  $scope.state = {};
  $scope.state.pagination_count = Pagination.getPaginationCount('swarm_nodes');
  $scope.sortType = 'Spec.Role';
  $scope.sortReverse = false;
  $scope.info = {};
  $scope.docker = {};
  $scope.swarm = {};
  $scope.totalCPU = 0;
  $scope.totalMemory = 0;

  $scope.order = function(sortType) {
    $scope.sortReverse = ($scope.sortType === sortType) ? !$scope.sortReverse : false;
    $scope.sortType = sortType;
  };

  $scope.changePaginationCount = function() {
    Pagination.setPaginationCount('swarm_nodes', $scope.state.pagination_count);
  };

  function extractSwarmInfo(info) {
    // Swarm info is available in SystemStatus object
    var systemStatus = info.SystemStatus;
    // Swarm strategy
    $scope.swarm[systemStatus[1][0]] = systemStatus[1][1];
    // Swarm filters
    $scope.swarm[systemStatus[2][0]] = systemStatus[2][1];
    // Swarm node count
    var nodes = systemStatus[0][1] === 'primary' ? systemStatus[3][1] : systemStatus[4][1];
    var node_count = parseInt(nodes, 10);
    $scope.swarm[systemStatus[3][0]] = node_count;

    $scope.swarm.Status = [];
    extractNodesInfo(systemStatus, node_count);
  }

  function extractNodesInfo(info, node_count) {
    // First information for node1 available at element #4 of SystemStatus if connected to a primary
    // If connected to a replica, information for node1 is available at element #5
    // The next 10 elements are information related to the node
    var node_offset = info[0][1] === 'primary' ? 4 : 5;
    for (i = 0; i < node_count; i++) {
      extractNodeInfo(info, node_offset);
      node_offset += 9;
    }
  }

  function extractNodeInfo(info, offset) {
    var node = {};
    node.name = info[offset][0];
    node.ip = info[offset][1];
    node.Id = info[offset + 1][1];
    node.status = info[offset + 2][1];
    node.containers = info[offset + 3][1];
    node.cpu = info[offset + 4][1].split('/')[1];
    node.memory = info[offset + 5][1].split('/')[1];
    node.labels = info[offset + 6][1];
    node.version = info[offset + 8][1];
    $scope.swarm.Status.push(node);
  }

  function processTotalCPUAndMemory(nodes) {
    var CPU = 0, memory = 0;
    angular.forEach(nodes, function(node) {
      CPU += node.CPUs;
      memory += node.Memory;
    });
    $scope.totalCPU = CPU / 1000000000;
    $scope.totalMemory = memory;
  }

  function initView() {
    $('#loadingViewSpinner').show();
  if (StateManager.getState().application.authentication) {
      var userDetails = Authentication.getUserDetails();
      var isAdmin = userDetails.role === 1 ? true: false;
      $scope.isAdmin = isAdmin;
  }
    var provider = $scope.applicationState.endpoint.mode.provider;
    $q.all({
      version: SystemService.version(),
      info: SystemService.info(),
      nodes: provider !== 'DOCKER_SWARM_MODE' || NodeService.nodes()
    })
    .then(function success(data) {
      $scope.docker = data.version;
      $scope.info = data.info;
      if (provider === 'DOCKER_SWARM_MODE') {
        var nodes = data.nodes;
        processTotalCPUAndMemory(nodes);
        $scope.nodes = nodes;
      } else {
        extractSwarmInfo(data.info);
      }
    })
    .catch(function error(err) {
      Notifications.error('Failure', err, 'Unable to retrieve cluster details');
    })
    .finally(function final() {
      $('#loadingViewSpinner').hide();
    });
  }

  initView();
}]);
