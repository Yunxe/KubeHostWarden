import React from "react";
import { Modal, Button } from "antd";

const AlertModal = ({ isVisible, handleOk, message }) => {
  return (
    <Modal
      title="操作结果"
      visible={isVisible}
      onOk={handleOk}
      onCancel={handleOk}
      footer={[
        <Button key="submit" type="primary" onClick={handleOk}>
          确认
        </Button>,
      ]}
    >
      <p>{message}</p>
    </Modal>
  );
};

export default AlertModal;
