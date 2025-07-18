"use client";

import { useRouter } from "next/navigation";
import SidebarButton from "./SidebarButton";
import { useState, useEffect, useRef } from "react";
import CreateWorkspacePopup from "./popup/CreateWorkspacePopup";
import CreateTeamPopup from "./popup/CreateTeamPopup";
import ChangeWorkspacePopup from "./popup/ChangeWorkspacePopup";
import RenameWorkspacePopup from "./popup/RenameWorkspacePopup";
import DeleteWorkspacePopup from "./popup/DeleteWorkspacePopup";

export default async function Sidebar(team: string) {
  const currentWorkspace = "Workspace 1";
  const owner = "John Doe";
  const teams = await fetch("http://localhost:8080/api/teams", {
    method: "GET",
  })
    .then((res) => res.json())
    .catch((err) => {
      console.error("Failed to fetch teams:", err);
    });
  const [selectedTeam, setSelectedTeam] = useState(team);
  const [popupNumber, setPopupNumber] = useState(0);
  const router = useRouter();
  const [showSidePopup, setShowSidePopup] = useState(false);
  const moreBtnRef = useRef<HTMLSpanElement>(null);
  const sidePopupRef = useRef<HTMLDivElement>(null);
  const [message, setMessage] = useState("");
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        sidePopupRef.current &&
        !sidePopupRef.current.contains(event.target as Node) &&
        moreBtnRef.current &&
        !moreBtnRef.current.contains(event.target as Node)
      ) {
        setShowSidePopup(false);
      }
    }
    if (showSidePopup) {
      document.addEventListener("mousedown", handleClickOutside);
    }
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, [showSidePopup]);
  const handleLogout = () => {
    router.push("/login");
  };
  const checkUser = (num: number) => {
    const user = localStorage.getItem("user");
    if (user !== owner) {
      alert("You do not have permission for this action.");
    } else setPopupNumber(num);
  };
  const selectTeam = (currentTeam: string) => {
    setSelectedTeam(currentTeam);
    router.push("/team");
  };
  const handlePopupSubmit = (value: string) => {
    console.log(value);
  };
  const handleCreateWorkspace = async (name: string) => {
    try {
      const res = await fetch("http://localhost:8080/api/workspace", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ name }),
      });

      if (!res.ok) {
        const errorText = await res.text();
        throw new Error(`API error: ${res.status} - ${errorText}`);
      }

      const data = await res.json();
      setMessage(data.message);
      alert(message);
      router.refresh();
    } catch (err) {
      console.error(err);
    }
  };
  const handleRenameWorkspace = async (name: string) => {
    const workspaceID = "workspace-id"; // Replace with actual workspace ID
    try {
      const res = await fetch("http://localhost:8080/api/workspace", {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ name }),
      });

      if (!res.ok) {
        const errorText = await res.text();
        throw new Error(`API error: ${res.status} - ${errorText}`);
      }

      const data = await res.json();
      setMessage(data.message);
      alert(message);
      router.refresh();
    } catch (err) {
      console.error(err);
    }
  };
  const handleDeleteWorkspace = async () => {
    const workspaceID = "workspace-id"; // Replace with actual workspace ID
    try {
      const res = await fetch("http://localhost:8080/api/workspace", {
        method: "DELETE",
      });

      if (!res.ok) {
        const errorText = await res.text();
        throw new Error(`API error: ${res.status} - ${errorText}`);
      }

      const data = await res.json();
      setMessage(data.message);
      alert(message);
      router.refresh();
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div className="flex-column items-center justify-between p-2 border-r-1 border-gray-500 bg-gray-800 text-gray-200 w-60 h-screen relative z-100">
      <div className="flex items-center justify-between py-2 mb-2">
        <span className="text-base font-semibold max-w-45 truncate">
          Username
        </span>
        <i
          className="fa-solid fa-right-from-bracket text-base"
          onClick={handleLogout}
        ></i>
      </div>
      <hr className="border-gray-500 mb-2"></hr>
      <CreateWorkspacePopup
        open={popupNumber === 1}
        onClose={() => setPopupNumber(0)}
        onSubmit={handleCreateWorkspace}
      />
      <ChangeWorkspacePopup
        open={popupNumber === 2}
        current={currentWorkspace}
        onClose={() => setPopupNumber(0)}
        onSubmit={handlePopupSubmit}
      />
      <CreateTeamPopup
        open={popupNumber === 3}
        onClose={() => setPopupNumber(0)}
        onSubmit={handlePopupSubmit}
      />
      <RenameWorkspacePopup
        open={popupNumber === 4}
        oldName={currentWorkspace}
        onClose={() => setPopupNumber(0)}
        onSubmit={handleRenameWorkspace}
      />
      <DeleteWorkspacePopup
        open={popupNumber === 5}
        name={currentWorkspace}
        onClose={() => setPopupNumber(0)}
        onSubmit={handleDeleteWorkspace}
      />
      <div>
        <div className="flex flex-row items-center justify-between p-1 mb-2 text-gray-400 text-sm">
          <span className="max-w-45 truncate">{currentWorkspace}</span>
          <span
            ref={moreBtnRef}
            className="w-6 h-6 text-center rounded-xl cursor-pointer hover:bg-gray-600 transition duration-200"
            onClick={() => setShowSidePopup(true)}
          >
            ...
          </span>
          {showSidePopup && (
            <div
              ref={sidePopupRef}
              className="absolute left-full top-15 bg-gray-700 text-gray-200 rounded shadow-lg p-2 w-45 z-50"
            >
              <div
                className="cursor-pointer hover:bg-gray-600 px-3 py-2 rounded transition duration-200"
                onClick={() => {
                  checkUser(4);
                  setShowSidePopup(false);
                }}
              >
                Rename workspace
              </div>
              <div
                className="cursor-pointer hover:bg-gray-600 px-3 py-2 rounded transition duration-200"
                onClick={() => {
                  checkUser(5);
                  setShowSidePopup(false);
                }}
              >
                Delete workspace
              </div>
            </div>
          )}
        </div>
        <div className="flex flex-col ml-4">
          <SidebarButton onClick={() => router.push("/workspace")}>
            Manage members
          </SidebarButton>
          <SidebarButton onClick={() => router.push("/my-issue")}>
            My issues
          </SidebarButton>
          <hr className="border-gray-500 mt-4 mb-4"></hr>
        </div>
      </div>
      <div className="block ml-4 overflow-y-auto max-h-115">
        <SidebarButton onClick={() => setPopupNumber(3)}>
          <i className="fa-solid fa-plus text-xs mr-2"></i>
          Create Team
        </SidebarButton>
        {teams.map((team: string[], index: number) => (
          <div key={index}>
            <SidebarButton
              onClick={() => selectTeam(team[1])}
              className={
                selectedTeam === team[1]
                  ? "bg-gray-700 text-white font-semibold"
                  : ""
              }
            >
              {team}
            </SidebarButton>
            {selectedTeam === team[1] && (
              <div className="ml-4">
                <SidebarButton
                  onClick={() => router.push("/project-list")}
                  className=""
                >
                  <i className="fa-solid fa-cube text-xs mr-4 w-2"></i>
                  Projects
                </SidebarButton>
                <SidebarButton
                  onClick={() => router.push("/issue-list")}
                  className=""
                >
                  <i className="fa-solid fa-bookmark text-xs mr-4 w-2"></i>
                  Issues
                </SidebarButton>
                <SidebarButton
                  onClick={() => router.push("/view-list")}
                  className=""
                >
                  <i className="fa-solid fa-layer-group text-xs mr-4 w-2"></i>
                  Views
                </SidebarButton>
              </div>
            )}
          </div>
        ))}
      </div>
      <div className="absolute bottom-0 left-0 w-full p-2 pb-6">
        <hr className="border-gray-500 mt-4 mb-4"></hr>
        <SidebarButton onClick={() => setPopupNumber(1)}>
          Create new workspace
        </SidebarButton>
        <SidebarButton onClick={() => setPopupNumber(2)}>
          Change workspace
        </SidebarButton>
      </div>
    </div>
  );
}
