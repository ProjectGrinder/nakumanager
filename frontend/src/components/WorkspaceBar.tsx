"use client";
export default function WorkspaceBar() {
  const workspaces = ["Workspace 1", "Workspace 2", "Workspace 3"];
  const joinedWorkspaces = ["Joined Workspace 1", "Joined Workspace 2"];
  const handleLogout = () => {
    // Logic for logging out the user
    console.log("User logged out");
  };
  const createWorkspace = () => {
    // Logic for creating a new workspace
    console.log("Create new workspace");
  };
  return (
    <div className="flex-column items-center justify-between p-4 border-r-1 border-gray-300 bg-gray-800 text-white w-60 h-screen">
      <div className="flex items-center justify-between p-2 pb-4 mb-4 border-b border-gray-600 text-xl font-bold">
        Username
        <i
          className="fa-solid fa-right-from-bracket text-xl"
          onClick={handleLogout}
        ></i>
      </div>
      <div className="flex items-center justify-between p-2 text-gray-200 text-sm font-normal">
        Your Workspaces
        <i
          className="fa-solid fa-square-plus text-xl"
          onClick={createWorkspace}
        ></i>
      </div>
      <ul className="flex-column items-start justify-start h-3/10 overflow-y-auto">
        {workspaces.map((workspace, index) => (
          <li key={index} className="p-2 hover:bg-gray-700 cursor-pointer">
            {workspace}
          </li>
        ))}
      </ul>
      <div className="flex items-center justify-between p-2 text-gray-200 text-sm font-normal">
        Joined Workspaces
      </div>
      <ul className="flex-column items-start justify-start h-1/2 overflow-y-auto">
        {joinedWorkspaces.map((joinedWorkspaces, index) => (
          <li key={index} className="p-2 hover:bg-gray-700 cursor-pointer">
            {joinedWorkspaces}
          </li>
        ))}
      </ul>
    </div>
  );
}
