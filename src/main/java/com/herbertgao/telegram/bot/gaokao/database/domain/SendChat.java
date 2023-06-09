package com.herbertgao.telegram.bot.gaokao.database.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import java.io.Serializable;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * 发送对话
 */
@Data
@AllArgsConstructor
@NoArgsConstructor
@TableName(value = "gaokao.send_chat")
public class SendChat implements Serializable {
    /**
     * ID
     */
    @TableId(value = "id", type = IdType.AUTO)
    private Long id;

    /**
     * 对话ID
     */
    @TableField(value = "chat_id")
    private String chatId;

    private static final long serialVersionUID = 1L;

    public static final String COL_ID = "id";

    public static final String COL_CHAT_ID = "chat_id";
}