package com.herbertgao.telegram.bot;

/**
 * 命令
 *
 * @author HerbertGao
 * @date 2019-06-08
 */
public interface Command {

    String COMMAND_REGEX = "^/[a-zA-Z]+";

    String COUNTDOWN_COMMAND = "/d";
    String LIST_COMMAND = "/ls";
    String ADD_COMMAND = "/add";
    String REMOVE_COMMAND = "/rm";
    String CUSTOMIZE_COMMAND = "/customize";
    String RENAME_COMMAND = "/rename";

}
