// import "./style.css";
import "./app.css";

// import logo from "./assets/images/logo-universal.png";
import { OnFileDrop } from "../wailsjs/runtime/runtime";
import {
  CreateImageDescriptor,
  CreateAndSavePortfolioDescriptorForImage,
  GetImagePreview,
  LoadDemoSeller,
  ReturnTagsForImage,
  SaveNewSellerTag,
  SearchFromImageDescriptor,
  SendBuyerChatMessage,
  SendSellerChatMessage,
} from "../wailsjs/go/main/App";

const sellerProfileId = 12; // Lawyer: 1 / Architect: 12
let allBuyerInput = "";

function setupUIEvents() {
  OnFileDrop(processFileDrop, true);

  const buyerInput = document.getElementById("buyer-input");
  buyerInput.addEventListener("keypress", async (e) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      const message = buyerInput.value.trim();

      if (message) {
        addMessage(message, "user");
        buyerInput.value = "";

        allBuyerInput += message;
        allBuyerInput += "\n";

        try {
          showThinking("buyer");

          SendBuyerChatMessage(allBuyerInput)
            .then((response) => {
              addMessage(response.followup, "bot");
              populateBuyerProfiles(response.profiles);
            })
            .catch((err) => {
              console.error(err);
            })
            .finally(() => {
              hideThinking("buyer");
            });
        } catch (err) {
          console.error(err);
          hideThinking("buyer");
        }
      }
    }
  });

  const tabButtons = document.querySelectorAll(".tab-button");
  tabButtons.forEach((button) => {
    button.addEventListener("click", () => {
      // Remove active class from all buttons and content
      tabButtons.forEach((btn) => btn.classList.remove("active"));
      document.querySelectorAll(".tab-content").forEach((content) => {
        content.classList.remove("active");
      });

      // Add active class to clicked button and corresponding content
      button.classList.add("active");
      const tabId = `${button.dataset.tab}-tab`;
      document.getElementById(tabId).classList.add("active");
    });
  });

  // Handle tag deletion
  document.querySelectorAll(".tag").forEach((tag) => {
    tag.addEventListener("click", function (e) {
      if (e.target.classList.contains("tag-delete")) {
        tag.remove();
      }
    });
  });

  // Handle adding new tags
  const sellerInput = document.querySelector(".seller-input");
  sellerInput.addEventListener("keypress", function (e) {
    if (e.key === "Enter" && this.value.trim()) {
      e.preventDefault();

      const message = this.value.trim();

      try {
        showThinking("seller");

        SendSellerChatMessage(message)
          .then((tags) => {
            for (const tag of tags) {
              addNewTag("seller-specialties-cloud", tag);
              SaveNewSellerTag(sellerProfileId, tag);
            }

            // Clear input
            this.value = "";
          })
          .catch((err) => {
            console.error(err);
          })
          .finally(() => {
            hideThinking("seller");
          });
      } catch (err) {
        console.error(err);
        hideThinking("seller");
      }
    }
  });
}

function setupUIData() {
  LoadDemoSeller(sellerProfileId).then((profile) => {
    populateSellerProfile(profile);
  });
}

document.addEventListener("DOMContentLoaded", function () {
  setupUIEvents();
  setupUIData();
});

function addMessage(text, type) {
  const chatHistory = document.getElementById("chat-history");
  const messageDiv = document.createElement("div");
  messageDiv.className = `message ${type}-message`;
  messageDiv.textContent = text;
  chatHistory.appendChild(messageDiv);
  chatHistory.scrollTop = chatHistory.scrollHeight;
}

function addNewTag(tagCloudDivId, tagValue) {
  const newTag = document.createElement("span");
  newTag.className = "tag";
  newTag.innerHTML = `${tagValue} <span class="tag-delete">×</span>`;

  newTag.addEventListener("click", function (e) {
    if (e.target.classList.contains("tag-delete")) {
      newTag.remove();
    }
  });

  document.getElementById(tagCloudDivId).appendChild(newTag);
}

function populateBuyerProfiles(profiles) {
  const profileSection = document.querySelector("#buyer-tab .profile-cards");
  profileSection.innerHTML = ""; // Clear existing profiles

  profiles.forEach((profile) => {
    const profileCard = document.createElement("div");
    profileCard.className = "profile-card";

    // Create a container for the SVG
    const imageContainer = document.createElement("div");
    imageContainer.innerHTML = profile.ProfileSVG;

    const profileInfo = document.createElement("div");
    profileInfo.className = "profile-info";

    // Create tags container
    const tagsContainer = document.createElement("div");
    tagsContainer.className = "tag-cloud";

    // Add specialty tags
    if (profile.SpecialtyTags && profile.SpecialtyTags.length > 0) {
      profile.SpecialtyTags.forEach((tag) => {
        const tagSpan = document.createElement("span");
        tagSpan.className = "tag";
        tagSpan.textContent = tag;
        tagsContainer.appendChild(tagSpan);
      });
    }

    profileInfo.innerHTML = `
        <div class="profile-title">${profile.Name}</div>
        <div class="vouches-counter">
            <span class="vouches-icon">👍</span>
            <span class="vouches-count">Vouched for by ${profile.Vouches || 0} Members in Your Network</span>
        </div>
        <div class="profile-summary">${profile.Summary}</div>
    `;

    // Add tags after the summary
    profileInfo.appendChild(tagsContainer);

    profileCard.appendChild(imageContainer.firstChild); // Add the SVG
    profileCard.appendChild(profileInfo);

    // Optional: Add click handler for profile cards
    profileCard.addEventListener("click", () => {
      handleProfileSelection(profile);
    });

    profileSection.appendChild(profileCard);
  });
}

function populateSellerProfile(profileData) {
  const profileImage = document.getElementById("seller-profile-image");
  const profileTitle = document.getElementById("seller-profile-title");
  const profileSummary = document.getElementById("seller-profile-summary");
  const specialtiesCloud = document.getElementById("seller-specialties-cloud");

  profileImage.outerHTML = profileData.ProfileSVG;
  profileTitle.textContent = profileData.Name;
  profileSummary.textContent = profileData.Summary;

  // Clear existing specialties
  specialtiesCloud.innerHTML = "";

  // Add new specialties
  profileData.SpecialtyTags.forEach((specialty) => {
    const tagSpan = document.createElement("span");
    tagSpan.className = "tag";

    const deleteSpan = document.createElement("span");
    deleteSpan.className = "tag-delete";
    deleteSpan.textContent = "×";

    // Add click handler for delete button
    deleteSpan.onclick = function () {
      tagSpan.remove();
    };

    tagSpan.textContent = specialty + " ";
    tagSpan.appendChild(deleteSpan);
    specialtiesCloud.appendChild(tagSpan);
  });
}

function processFileDrop(x, y, paths) {
  // only process first file
  const path = paths[0];

  const isImage = path.match(/\.(jpg|jpeg|png|gif|bmp)$/i);

  const isBuyerActive = document
    .getElementById("buyer-tab")
    .classList.contains("active");
  const isSellerActive = document
    .getElementById("seller-tab")
    .classList.contains("active");

  const previewImage = isBuyerActive
    ? document.getElementById("preview-image-buyer")
    : document.getElementById("preview-image-seller");

  if (isImage) {
    GetImagePreview(path)
      .then((base64String) => {
        previewImage.style.display = "block";
        previewImage.src = `data:image/jpeg;base64,${base64String}`;
      })
      .catch((err) => {
        console.error("Error loading preview:", err);
      });
  }

  if (isBuyerActive) {
    showThinking("buyer");

    SearchFromImageDescriptor(path)
      .then((profiles) => {
        populateBuyerProfiles(profiles);
      })
      .catch((err) => {
        console.error(err);
      })
      .finally(() => {
        hideThinking("buyer");
      });
  } else if (isSellerActive) {
    showThinking("seller");

    CreateAndSavePortfolioDescriptorForImage(sellerProfileId, path)
      .then(() => {})
      .catch((err) => {
        console.error(err);
      })
      .finally(() => {
        hideThinking("seller");
      });
  }
}

function showThinking(section = "buyer") {
  const selector =
    section === "buyer"
      ? "#chat-history .thinking-indicator"
      : ".seller-input-column .thinking-indicator";
  const indicator = document.querySelector(selector);
  indicator.classList.remove("hidden");
}

function hideThinking(section = "buyer") {
  const selector =
    section === "buyer"
      ? "#chat-history .thinking-indicator"
      : ".seller-input-column .thinking-indicator";
  const indicator = document.querySelector(selector);
  indicator.classList.add("hidden");
}
