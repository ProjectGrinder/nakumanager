"use client";

import { useRouter } from "next/navigation";
import SidebarButton from "./SidebarButton";
import { useState, useEffect, useRef } from "react";
import CreateWorkspacePopup from "./popup/CreateWorkspacePopup";
import CreateTeamPopup from "./popup/CreateTeamPopup";
import ChangeWorkspacePopup from "./popup/ChangeWorkspacePopup";
import RenameWorkspacePopup from "./popup/RenameWorkspacePopup";

export default function Sidebar(team: string) {
  const currentWorkspace = "Workspace 1";
  const teams = ["Team 1", "Team 2", "Team 3"];
  const [selectedTeam, setSelectedTeam] = useState(team);
  const [popupNumber, setPopupNumber] = useState(0);
  const router = useRouter();
  const [showSidePopup, setShowSidePopup] = useState(false);
  const moreBtnRef = useRef<HTMLSpanElement>(null);
  const sidePopupRef = useRef<HTMLDivElement>(null);
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
  const selectTeam = (currentTeam: string) => {
    setSelectedTeam(currentTeam);
    router.push("/team");
  };
  const handlePopupSubmit = (value: string) => {
    console.log(value);
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
        onSubmit={handlePopupSubmit}
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
        onSubmit={handlePopupSubmit}
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
                  setPopupNumber(4);
                  setShowSidePopup(false);
                }}
              >
                Rename workspace
              </div>
              <div
                className="cursor-pointer hover:bg-gray-600 px-3 py-2 rounded transition duration-200"
                onClick={() => {
                  setPopupNumber(5);
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
      <div className="flex flex-col ml-4">
        <SidebarButton onClick={() => setPopupNumber(3)}>
          <i className="fa-solid fa-plus text-xs mr-2"></i>
          Create Team
        </SidebarButton>
        {teams.map((team, index) => (
          <div key={index}>
            <SidebarButton
              onClick={() => selectTeam(team)}
              className={
                selectedTeam === team
                  ? "bg-gray-700 text-white font-semibold"
                  : ""
              }
            >
              {team}
            </SidebarButton>
            {selectedTeam === team && (
              <div className="flex flex-col ml-4">
                <SidebarButton onClick={() => router.push("/project-list")}>
                  <i className="fa-solid fa-cube text-xs mr-4 w-2"></i>
                  Projects
                </SidebarButton>
                <SidebarButton onClick={() => router.push("issue-list")}>
                  <i className="fa-solid fa-bookmark text-xs mr-4 w-2"></i>
                  Issues
                </SidebarButton>
                <SidebarButton onClick={() => router.push("view-list")}>
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
