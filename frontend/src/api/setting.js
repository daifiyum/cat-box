import service from "@/utils/axios";

// 获取配置
export function getSetting(key) {
  return service({
    method: "GET",
    url: `/setting/${key}`,
  });
}

// 更新配置
export function setSetting(key, value) {
  return service({
    method: "POST",
    url: `/setting/${key}`,
    data: { value },
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
  });
}
